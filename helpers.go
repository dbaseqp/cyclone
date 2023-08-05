package main

import (
	// "errors"
	// "log"
	"math/rand"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyz")
	uppers  = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	symbols = []rune("!@#$")
)

func GeneratePassword(n int) string {
	rand.Seed(time.Now().UnixNano())

	var b []rune
	for i := 0; i < n/2; i++ {
		b = append(b, letters[rand.Intn(len(letters))])
	}
	for i := 0; i < n/2; i++ {
		b = append(b, uppers[rand.Intn(len(uppers))])
	}
	b = append(b, symbols[rand.Intn(len(symbols))])
	return string(b)
}

func GenerateToken(username string) (string, error) {
	token_lifespan := 1

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tomlConf.JwtSecret))
}

func ValidateToken(c *gin.Context) error {
	tokenString := ExtractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tomlConf.JwtSecret), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenUsername(c *gin.Context) (string, error) {

	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tomlConf.JwtSecret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims["username"], nil
	}
	return "", nil
}