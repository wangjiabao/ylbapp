package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	UserId   int64
	UserType string
	jwt.StandardClaims
}

// CreateToken generate token
func CreateToken(c CustomClaims, key string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signedString, err := claims.SignedString([]byte(key))
	if err != nil {
		return "", errors.New("generate token failed" + err.Error())
	}
	return signedString, nil
}
