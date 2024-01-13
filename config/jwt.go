package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("Rahasia30193103910admnakdnak1391039")

type JWTClaim struct {
	Username string
	jwt.RegisteredClaims
}
