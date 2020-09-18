package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"

	C "github.com/sumaikun/apeslogistic-rest-api/config"

	Dao "github.com/sumaikun/apeslogistic-rest-api/dao"

	"github.com/gorilla/context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var dao = Dao.MongoConnector{}

// AuthMiddleware verify
func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {		

		at := r.Header.Get("Authorization")

		ifContains := strings.Contains(at, "Bearer ")

		if ifContains == true {

			cognitoClient := cognito.GetUserInput{ AccessToken: at }

			err := cognitoClient.Validate()

			if err != nil {
				fmt.Printf(err)
				return
			}

			next.ServeHTTP(w, r)
			return

		}

	
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
		return jwtMiddleware.Handler(next)

	}	

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

		if userParsed["role"] == "admin" {
			next.ServeHTTP(w, r)
		} else {
			return
		}

	})

}
