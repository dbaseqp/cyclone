package main

import (
	"errors"
	// "log"
	"math/rand"
	"time"
	"fmt"
	"strings"
	"context"
	"net/url"
	"os/exec"
	"bytes"
	"regexp"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
		return claims["username"].(string), nil
	}
	return "", nil
}

func ValidateVsphereCredentials(username string, password string) error {
	ctx := context.Background()
	u, err := soap.ParseURL(tomlConf.VCenterURL)
	if err != nil {
		fmt.Printf("Error parsing vCenter URL: %s\n", err)
		return err
	}

	u.User = url.UserPassword(username, password)

	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		fmt.Printf("Error creating vSphere client for %s: %s\n", username, err)
		return err
	}
	defer client.Logout(ctx)
	return nil
}

func ValidateJWT(token string) error {
	ctx := context.Background()
	u, err := soap.ParseURL(tomlConf.VCenterURL)
	if err != nil {
		fmt.Printf("Error parsing vCenter URL: %s\n", err)
		return err
	}
	creds := strings.Split(token, ":")
	username := creds[0]
	password := creds[1]
	u.User = url.UserPassword(username, password)

	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		fmt.Printf("Error creating vSphere client for %s: %s\n", username, err)
		return err
	}
	defer client.Logout(ctx)
	return nil
}

func RegisterUser(username string, password string) error {
	matched, _ := regexp.MatchString(`^\w{1,16}$`, username)

	if !matched {
		return errors.New("Username must not exceed 16 characters and may only contain letters, numbers, or an underscore (_)!")
	}

	cmd := exec.Command("powershell", "New-PodUser", username, fmt.Sprintf("'%s'", password))

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return errors.New("Failed to register user!")
	}
	
	return nil
}