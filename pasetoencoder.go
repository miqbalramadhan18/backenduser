package backenduser

import (
	"encoding/json"

	"fmt"

	"aidanwoods.dev/go-paseto"
)

// <--- paseto encoder --->
// memeriksa token
func IsTokenValid(publickey, tokenstr string) (payload Payload, err error) {
	var token *paseto.Token
	var pubKey paseto.V4AsymmetricPublicKey
	pubKey, err = paseto.NewV4AsymmetricPublicKeyFromHex(publickey) // this wil fail if given key in an invalid format
	if err != nil {
		fmt.Println("Decode NewV4AsymmetricPublicKeyFromHex : ", err)
	}
	parser := paseto.NewParser()                             // only used because this example token has expired, use NewParser() (which checks expiry by default)
	token, err = parser.ParseV4Public(pubKey, tokenstr, nil) // this will fail if parsing failes, cryptographic checks fail, or validation rules fail
	if err != nil {
		fmt.Println("Decode ParseV4Public : ", err)
	} else {
		json.Unmarshal(token.ClaimsJSON(), &payload)
	}
	return payload, err
}
