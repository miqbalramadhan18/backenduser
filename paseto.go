package backenduser

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/whatsauth/watoken"
)

// <--- ini Login & Register User --->
func Login(Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp Credential
	mconn := SetConnection(MongoEnv, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, Colname, datauser) {
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(Privatekey))
			if err != nil {
				resp.Message = "Gagal Encode Token : " + err.Error()
			} else {
				resp.Status = true
				resp.Message = "Selamat Datang User"
				resp.Token = tokenstring
			}
		} else {
			resp.Message = "Password Salah"
		}
	}
	return GCFReturnStruct(resp)
}


// return struct
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func ReturnStringStruct(Data any) string {
	jsonee, _ := json.Marshal(Data)
	return string(jsonee)
}

func Register(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	userdata := new(User)
	resp.Status = false
	conn := SetConnection(Mongoenv, dbname)
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		hash, err := HashPass(userdata.Password)
		if err != nil {
			resp.Message = "Gagal Hash Password" + err.Error()
		}
		InsertUserdata(conn, userdata.Username, userdata.Role, hash)
		resp.Message = "Berhasil Input data"
	}
	response := ReturnStringStruct(resp)
	return response
}


// <--- ini Report --->

//Create Report post
// func GCFCreateReport(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var datareport Report
// 	err := json.NewDecoder(r.Body).Decode(&datareport)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	if err := CreateReport(mconn, collectionname, datareport); err != nil {
// 		return GCFReturnStruct(CreateResponse(true, "Success Create Catalog", datareport))
// 	} else {
// 		return GCFReturnStruct(CreateResponse(false, "Failed Create Catalog", datareport))
// 	}
// }

// Insert Report post 
func GCFInsertReport(publickey, MONGOCONNSTRINGENV, dbname, colluser, collreport string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var userdata User
	gettoken := r.Header.Get("token")
	if gettoken == "" {
		response.Message = "Missing token in headers"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		userdata.Username = checktoken
		if checktoken == "" {
			response.Message = "Invalid token"
		} else {
			user2 := FindUser(mconn, colluser, userdata)
			if user2.Role == "user" {
				var datareport Report
				err := json.NewDecoder(r.Body).Decode(&datareport)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					insertReport(mconn, collreport, Report{
						Nik:     		datareport.Nik,
						Title:       	datareport.Title,
						Description: 	datareport.Description,
						DateOccurred: 	datareport.DateOccurred,
						Image:       	datareport.Image,
						Status:      	datareport.Status,
					})
					response.Status = true
					response.Message = "Berhasil Insert Report"
				}
			} else {
				response.Message = "Anda tidak bisa Insert data karena bukan User"
			}
		}
	}
	return GCFReturnStruct(response)
}

// delete report
func GCFDeleteReport(publickey, MONGOCONNSTRINGENV, dbname, colluser, collreport string, r *http.Request) string {

	var respon Credential
	respon.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var userdata User

	gettoken := r.Header.Get("token")
	if gettoken == "" {
		respon.Message = "Missing token in headers"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		userdata.Username = checktoken
		if checktoken == "" {
			respon.Message = "Invalid token"
		} else {
			user2 := FindUser(mconn, colluser, userdata)
			if user2.Role == "user" {
				var datareport Report
				err := json.NewDecoder(r.Body).Decode(&datareport)
				if err != nil {
					respon.Message = "Error parsing application/json: " + err.Error()
				} else {
					DeleteReport(mconn, collreport, datareport)
					respon.Status = true
					respon.Message = "Berhasil Delete Report"
				}
			} else {
				respon.Message = "Anda tidak bisa Delete data karena bukan User"
			}
		}
	}
	return GCFReturnStruct(respon)
}

// // update Report
// func GCFUpdateReport(publickey, MONGOCONNSTRINGENV, dbname, colluser, collreport string, r *http.Request) string {
// 	var response Credential
// 	response.Status = false
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var userdata User

// 	gettoken := r.Header.Get("token")
// 	if gettoken == "" {
// 		response.Message = "Missing token in Headers"
// 	} else {
// 		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
// 		userdata.Username = checktoken
// 		if checktoken == "" {
// 			response.Message = "Invalid token"
// 		} else {
// 			user2 := FindUser(mconn, colluser, userdata)
// 			if user2.Role == "user" {
// 				var datareport Report
// 				err := json.NewDecoder(r.Body).Decode(&datareport)
// 				if err != nil {
// 					response.Message = "Error parsing application/json: " + err.Error()

// 				} else {
// 					UpdatedReport(mconn, collreport, bson.M{"id": datareport.ID}, datareport)
// 					response.Status = true
// 					response.Message = "Berhasil Update Report"
// 					GCFReturnStruct(CreateResponse(true, "Success Update Report", datareport))
// 				}
// 			} else {
// 				response.Message = "Anda tidak bisa Update data karena bukan User"
// 			}

// 		}
// 	}
// 	return GCFReturnStruct(response)
// }


// get all report
func GCFGetAllReport(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datareport := GetAllReport(mconn, collectionname)
	if datareport != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Report", datareport))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Report", datareport))
	}
}

// func GetAllDataReports(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
// 	req := new(Response)
// 	conn := SetConnection(MongoEnv, dbname)
// 	tokenlogin := r.Header.Get("Login")
// 	if tokenlogin == "" {
// 		req.Status = false
// 		req.Message = "Header Login Not Found"
// 	} else {
// 		// Dekode token untuk mendapatkan
// 		_, err := DecodeGetReport(os.Getenv(PublicKey), tokenlogin)
// 		if err != nil {
// 			req.Status = false
// 			req.Message = "Tidak ada data  " + tokenlogin
// 		} else {
// 			// Langsung ambil data report
// 			datareport := GetAllReport(conn, colname)
// 			if datareport == nil {
// 				req.Status = false
// 				req.Message = "Data Report tidak ada"
// 			} else {
// 				req.Status = true
// 				req.Message = "Data Report berhasil diambil"
// 				req.Data = datareport
// 			}
// 		}
// 	}
// 	return ReturnStringStruct(req)
// }

// get all report by id
func GCFGetAllReportID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datareport Report
	err := json.NewDecoder(r.Body).Decode(&datareport)
	if err != nil {
		return err.Error()
	}

	report := GetAllReportID(mconn, collectionname, datareport)
	if report != (Report{}) {
		return GCFReturnStruct(CreateResponse(true, "Success: Get ID Report", datareport))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed to Get ID Report", datareport))
	}
}
