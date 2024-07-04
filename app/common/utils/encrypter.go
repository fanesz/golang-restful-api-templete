package utils

import "golang.org/x/crypto/bcrypt"

func Encrypt(text *string) {
	res, err := bcrypt.GenerateFromPassword([]byte(*text), bcrypt.DefaultCost)
	if err != nil {
		*text = ""
	}
	*text = string(res)
}

func CompareEncrypted(text *string, textToCompare *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*text), []byte(*textToCompare))
	return err == nil
}
