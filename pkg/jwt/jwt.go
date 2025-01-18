package jwt

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(userID int) string {
	expirationDate := time.Now().Add(24 * time.Hour).UnixMilli()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":       "knights-vow",
		"sub":       "access",
		"aud":       "knights-vow-client",
		"exp":       expirationDate,
		"client_id": userID,
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		log.Fatalf("Error signing token: %v", err)
	}

	return tokenString
}

func ParseJWT(tokenString string) *jwt.Token {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, err := token.Method.(*jwt.SigningMethodHMAC); !err {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if err != nil {
		log.Fatalf("Error parsing token: %v", err)
	}

	return token
}

func Verify(tokenString string) bool {
	token := ParseJWT(tokenString)

	return token.Valid
}
