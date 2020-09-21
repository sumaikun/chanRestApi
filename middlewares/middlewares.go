package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"

	C "github.com/sumaikun/apeslogistic-rest-api/config"

	Dao "github.com/sumaikun/apeslogistic-rest-api/dao"

	"github.com/gorilla/context"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	Helpers "github.com/sumaikun/apeslogistic-rest-api/helpers"
)

var dao = Dao.MongoConnector{}

// AuthMiddleware verify
func AuthMiddleware(next http.Handler) http.Handler {

	fmt.Println()

	var config = C.Config{}
	config.Read()

	var JwtKey = []byte(config.Jwtkey)

	if len(JwtKey) == 0 {
		log.Fatal("HTTP server unable to start, expected an APP_KEY for JWT auth")
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		Extractor: jwtmiddleware.FromFirst(jwtmiddleware.FromAuthHeader,
			jwtmiddleware.FromParameter("token")),
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cognito_email := context.Get(r, "cognito_email")

		fmt.Printf("cognito_email", cognito_email)

		if cognito_email == nil {
			err := jwtMiddleware.CheckJWT(w, r)
			if err != nil {
				Helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
		return
	})

	//return jwtMiddleware.Handler(next)

}

// CognitoMiddleware verify
func CognitoMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		at := r.Header.Get("Authorization")

		fmt.Println("at", at)

		ifContains := strings.Contains(at, "Bearer ")

		fmt.Println("ifContains", ifContains)

		if ifContains == false {

			fmt.Println("lets check cognito")

			conf := &aws.Config{Region: aws.String("us-east-1")}

			sess, err := session.NewSession(conf)
			if err != nil {
				Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			cognitoClient := cognito.New(sess)

			userInput := cognito.GetUserInput{AccessToken: &at}

			cognitoResponse, err2 := cognitoClient.GetUser(&userInput)

			if err2 != nil {
				Helpers.RespondWithError(w, http.StatusUnauthorized, err2.Error())
				return
			}

			fmt.Println("cognitoResponse", cognitoResponse.UserAttributes)

			for _, s := range cognitoResponse.UserAttributes {
				//fmt.Println(i, s)
				if *s.Name == "email" {
					context.Set(r, "cognito_email", s.Value)
				}
			}

			next.ServeHTTP(w, r)
			return

		}

		next.ServeHTTP(w, r)
		return

	})

}

// UserMiddleware get user from request
func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var config = C.Config{}
		config.Read()

		ua := r.Header.Get("Authorization")

		ua = strings.Replace(ua, "Bearer ", "", 1)

		tokenString := ua
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Jwtkey), nil
		})
		// ... error handling

		if err != nil {
			log.Fatal("Error decoding jwt")
		}

		//log.Println(claims["username"])

		user, err := dao.FindOneByKEY("users", "email", claims["username"].(string))

		if err != nil {
			log.Fatal("Can not get user from token")
			return
		}

		context.Set(r, "user", user)

		//log.Println(user)

		next.ServeHTTP(w, r)

		//log.Println("Executing middlewareOne again")
	})
}

// OnlyAdminMiddleware can execute request if is admin
func OnlyAdminMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := context.Get(r, "user")

		userParsed := user.(bson.M)

		if userParsed["role"] == "ADMIN" {
			next.ServeHTTP(w, r)
		} else {
			Helpers.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "only admin can do this work"})
			return
		}

	})

}
