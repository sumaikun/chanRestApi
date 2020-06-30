package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	Config "github.com/sumaikun/apeslogistic-rest-api/config"
)

var (
	port   string
	jwtKey []byte
)

func init() {
	var config = Config.Config{}
	config.Read()
	//fmt.Println(config.Jwtkey)
	jwtKey = []byte(config.Jwtkey)
	port = config.Port
}

func main() {

	fmt.Println("start server in port " + port)
	router := mux.NewRouter().StrictSlash(true)

	/* Authentication */
	router.HandleFunc("/auth", authentication).Methods("POST")

	/* Participants */
	router.HandleFunc("/participants", authentication).Methods("GET")

	/* Assets */
	router.HandleFunc("/assets", authentication).Methods("GET")

	/* ISSUES */
	router.HandleFunc("/issues", authentication).Methods("GET")

	/* TRAZABILITY */
	//router.HandleFunc("/traz/{id}", authentication).Methods("GET")

	//start server
	log.Fatal(http.ListenAndServe(":"+port, router))

}
