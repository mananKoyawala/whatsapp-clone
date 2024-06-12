package helper

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

var SECRET_KEY string

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("can't load env")
	}

	SECRET_KEY = os.Getenv("SECRET_KEY")
}

type SignedDetails struct {
	ID int64
	jwt.StandardClaims
}

func GenerateJwtToken(id int64) (string, string, error) {

	// Create the Claims
	claims := &SignedDetails{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refershClaims := &SignedDetails{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refershClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	// log.Println(SECRET_KEY)
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return nil, err.Error()
	}

	//Token is invalid
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The token is invalid"
		return claims, msg
	}

	// Token is expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is expired "
		return claims, msg
	}

	return claims, ""
}
