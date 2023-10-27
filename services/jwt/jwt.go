package jwt

import (
	"escrolla-api/errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"time"
)

const AccessTokenValidity = time.Minute * 5
const PasswordReset = time.Minute * 20
const RefreshTokenValidity = time.Minute * 30
const ConfirmEmailValidity = time.Hour * 24 * 2

func getClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not get claims")
	}
	return claims, claims.Valid()
}

func ValidateAndGetClaims(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := validateToken(tokenString, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	claims, err := getClaims(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get claims: %v", err)
	}
	return claims, nil
}

func GetTokenFromHeader(req *http.Request) string {
	authHeader := req.Header.Get("Authorization")
	if len(authHeader) > 8 {
		return authHeader[7:]
	}
	return ""
}

func GenerateToken(email, secret string, expiryDuration time.Duration) (string, error) {
	if secret == "" {
		return "", errors.New("", http.StatusInternalServerError)
	}
	// Generate claims
	claims := generateClaims(email, expiryDuration)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// verifyAccessToken verifies a token
func validateToken(tokenString string, secret string) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, errors.New("invalid token (token is empty)", http.StatusForbidden)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token", http.StatusForbidden)
	}
	return token, err
}

func generateClaims(email string, expiryDuration time.Duration) jwt.MapClaims {
	accessClaims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(expiryDuration).Unix(),
	}
	return accessClaims
}
