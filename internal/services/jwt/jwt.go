package jwtservice

import (
	"fmt"
	"oauth-go/internal/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	IssueTokensPair() (string, string)
	VerifyToken() (*jwt.Token, error)
}

type Jwt struct {
	config *types.AppConfig
}

func New(config *types.AppConfig) *Jwt {
	return &Jwt{
		config: config,
	}
}

type AppCustomClaims struct {
	UserID    int    `json:"user_id"`
	Email     string `json:"email"`
	SessionID int    `json:"session_id"`
}

type CustomClaims struct {
	AppCustomClaims
	jwt.RegisteredClaims
}

func (service *Jwt) IssueTokensPair(userID int, sessionID int, email string) (string, string) {
	// Create access token (short-lived JWT)
	accessTokenClaims := &CustomClaims{
		AppCustomClaims: AppCustomClaims{
			UserID:    userID,
			SessionID: sessionID,
			Email:     email,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	secret := []byte(service.config.JwtSecret)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, _ := accessToken.SignedString(secret)

	// Create refresh token (long-lived JWT)
	refreshTokenClaims := &CustomClaims{
		AppCustomClaims: AppCustomClaims{
			UserID:    userID,
			SessionID: sessionID,
			Email:     email,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 7)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, _ := refreshToken.SignedString(secret)

	return accessTokenString, refreshTokenString
}

func (service *Jwt) VerifyToken(tokenString string) (*jwt.Token, error) {
	secret := []byte(service.config.JwtSecret)
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		// Ensure the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Check if the token is expired
	if claims, ok := token.Claims.(CustomClaims); ok {
		exp, err := claims.GetExpirationTime()

		if err != nil {
			return nil, err
		}

		if time.Now().After(exp.UTC()) {
			return nil, fmt.Errorf("token has expired")
		}
	}

	return token, nil
}

func (service *Jwt) GetClaims(token *jwt.Token) (*CustomClaims, error) {
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("cannot cast token.Claims to *CustomClaims")
	}

	return claims, nil
}
