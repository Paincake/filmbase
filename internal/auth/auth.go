package auth

import (
	"github.com/golang-jwt/jwt"
	"os"
)

func CreateJWT(username, role string) (string, error) {
	var sampleSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["role"] = role
	claims["username"] = username
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
