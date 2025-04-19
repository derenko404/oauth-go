package services

import (
	jwtservice "oauth-go/internal/services/jwt"
	ouathservice "oauth-go/internal/services/oauth"
	"oauth-go/internal/types"
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
