package model

import "github.com/dgrijalva/jwt-go"

type UserJwt struct {
	UserId int64  `json:"userId"`
	Role   string `json:"role"`
}

type UserClaims struct {
	jwt.StandardClaims
	UserId int64  `json:"userId"`
	Role   string `json:"role"`
}
