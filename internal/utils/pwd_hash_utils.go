package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// Hashing the password with a cost of 14 (adjust as needed)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CompareHashPwd(hashpwd string, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashpwd), []byte(pwd))
}
