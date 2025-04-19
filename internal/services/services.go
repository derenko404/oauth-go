package services

import (
	jwtservice "go-auth/internal/services/jwt"
	ouathservice "go-auth/internal/services/oauth"
	"go-auth/internal/types"
)

type Services struct {
	OAuth *ouathservice.OAuth
	Jwt   *jwtservice.Jwt
}

func New(config *types.AppConfig) *Services {
	return &Services{
		OAuth: ouathservice.New(config),
		Jwt:   jwtservice.New(config),
	}
}
