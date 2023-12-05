package backenduser

import (
	"context"
	"fmt"
	"os"

	"github.com/aiteung/atdb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetConnection(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		// DBString: "mongodb+srv://miqbalr18:18112002@qobel.loal6fn.mongodb.net/", //os.Getenv(MONGOCONNSTRINGENV),
		DBString: os.Getenv(MONGOCONNSTRINGENV),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

func InsertUserdata(MongoConn *mongo.Database, username, role, password string) (InsertedID interface{}) {
	req := new(User)
	req.Username = username
	req.Password = password
	req.Role = role
	return InsertOneDoc(MongoConn, "user", req)
}

func DeleteUser(mongoconn *mongo.Database, collection string, userdata User) interface{} {
	filter := bson.M{"username": userdata.Username}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func FindUser(mongoconn *mongo.Database, collection string, userdata User) User {
	filter := bson.M{"username": userdata.Username}
	return atdb.GetOneDoc[User](mongoconn, collection, filter)
}

func IsPasswordValid(mongoconn *mongo.Database, collection string, userdata User) bool {
	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mongoconn, collection, filter)
	return CompareHashPass(userdata.Password, res.Password)
}

func MongoCreateConnection(MongoString, dbname string) *mongo.Database {
	MongoInfo := atdb.DBInfo{
		DBString: os.Getenv(MongoString),
		DBName:   dbname,
	}
	conn := atdb.MongoConnect(MongoInfo)
	return conn
}

func InsertOneDoc(db *mongo.Database, collection string, doc interface{}) (insertedID interface{}) {
	insertResult, err := db.Collection(collection).InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	return insertResult.InsertedID
}

func GetOneUser(MongoConn *mongo.Database, colname string, userdata User) User {
	filter := bson.M{"username": userdata.Username}
	data := atdb.GetOneDoc[User](MongoConn, colname, filter)
	return data
}

