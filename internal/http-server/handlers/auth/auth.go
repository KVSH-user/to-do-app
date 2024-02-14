package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
)

var mySigningKey = "dhaw7dyaw8"

func ValidateToken(tokenString string) (string, *jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(mySigningKey), nil // Используйте []byte для ключа
	})

	if err != nil {
		return "", token, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid, ok := claims["uid"].(float64) // UID может быть float64 из-за спецификации JSON
		if !ok {
			return "", token, fmt.Errorf("uid claim is missing")
		}
		return strconv.Itoa(int(uid)), token, nil // Преобразование UID обратно в строку
	} else {
		return "", token, fmt.Errorf("invalid token")
	}
}
