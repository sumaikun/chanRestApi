package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	Helpers "github.com/sumaikun/apeslogistic-rest-api/helpers"
	Models "github.com/sumaikun/apeslogistic-rest-api/models"
)

func authentication(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var creds Models.Credentials

	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := Models.Users[creds.Username]

	//fmt.Println(creds)

	if !ok || !Helpers.CheckPasswordHash(creds.Password, expectedPassword) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(8 * time.Hour)

	claims := &Models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")

	//Generate json response for get the token
	json.NewEncoder(w).Encode(&Models.TokenResponse{Token: tokenString})

}
