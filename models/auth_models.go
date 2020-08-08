package models

import (
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

// Credentials is the request body of credential input
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Users are test users for generate jwt token
var Users = map[string]string{
	"sumaikun":                  "$2a$14$6NWsioRmg3dogylbm0j3e.0RVDAN2dybn2HzecrFCNex9PPxsEJLi",
	"user2":                     "$2a$14$6NWsioRmg3dogylbm0j3e.0RVDAN2dybn2HzecrFCNex9PPxsEJLi",
	"jesus.vega@technoapes.com": "$2a$14$csInuvCwTayoNABvSvxH6uPdDF11jS/Uh2I/nrEtTDbZC3cpL5nDy",
}

// Claims represents the struct of jwt token
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// TokenResponse represents json response after succesfully auth
type TokenResponse struct {
	Token string `json:"token"`
	User  bson.M `json:"user"`
}
