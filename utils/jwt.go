package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtSecret []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set in .env")
	}
	jwtSecret = []byte(secret)
}

type JWTClaims struct {
	UserID     int    `json:"user_id"`
	EmployeeID int    `json:"employee_id"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token valid for 120 hours
func GenerateToken(userID int, employeeID int, role string) (string, error) {
	claims := JWTClaims{
		UserID:     userID,
		EmployeeID: employeeID,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(120 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
