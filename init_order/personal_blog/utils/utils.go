package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"personal_blog/model"
	"time"
)

func HashPassWord(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func ValidPassWord(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(user model.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.UserName,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	return claims.SignedString([]byte("secret"))
}
func ParseToken(token string) ([]interface{}, error) {
	parse, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("the token type is incorrect")
		} else {
			return []byte("secret"), nil
		}
	})
	if err != nil {
		return nil, err
	}
	if claim, ok := parse.Claims.(jwt.MapClaims); ok && parse.Valid {
		username, ok := claim["username"].(string)
		//fmt.Println(claim["id"])
		id, ok1 := claim["id"].(float64)
		if !ok {
			return nil, errors.New("username claims is not string")
		}
		if !ok1 {
			return nil, errors.New("id claims is not float64")
		}
		return []interface{}{id, username}, nil
	}
	return nil, err
}
