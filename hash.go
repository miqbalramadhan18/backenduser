package backenduser

import "golang.org/x/crypto/bcrypt"

func HashPass(password string) (string, error) {
	bytess, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytess), err
}

func CompareHashPass(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
