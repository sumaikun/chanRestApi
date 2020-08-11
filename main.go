package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	//"github.com/chainHero/heroes-service/blockchain"
	"github.com/gorilla/mux"
	"github.com/sumaikun/apeslogistic-rest-api/blockchain"
	Config "github.com/sumaikun/apeslogistic-rest-api/config"
	Dao "github.com/sumaikun/apeslogistic-rest-api/dao"
	middleware "github.com/sumaikun/apeslogistic-rest-api/middlewares"
)

//Application object to chaincode connection
type Application struct {
	Fabric *blockchain.FabricSetup
}

var (
	port   string
	jwtKey []byte
)

var dao = Dao.MongoConnector{}

var fSetup = blockchain.FabricSetup{
	// Network parameters
	OrdererID: "orderer.hf.chainhero.io",

	// Channel parameters
	ChannelID:     "chainhero",
	ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/chainHero/heroes-service/fixtures/artifacts/chainhero.channel.tx",

	// Chaincode parameters
	ChainCodeID:     "heroes-service",
	ChaincodeGoPath: os.Getenv("GOPATH"),
	ChaincodePath:   "github.com/sumaikun/apeslogistic-rest-api/chaincode/",
	OrgAdmin:        "Admin",
	OrgName:         "org1",
	ConfigFile:      "config.yaml",

	// User parameters
	UserName: "User1",
}

func init() {

	var config = Config.Config{}
	config.Read()
	//fmt.Println(config.Jwtkey)
	jwtKey = []byte(config.Jwtkey)
	port = config.Port

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()

}

// CORSRouterDecorator applies CORS headers to a mux.Router
type CORSRouterDecorator struct {
	R *mux.Router
}

// ServeHTTP wraps the HTTP server enabling CORS headers.
// For more info about CORS, visit https://www.w3.org/TR/cors/
func (c *CORSRouterDecorator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	//fmt.Println("I am on serve HTTP")

	rw.Header().Set("Access-Control-Allow-Origin", "*")

	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		//fmt.Println("I am in options")
		rw.WriteHeader(http.StatusOK)
		return
	}

	c.R.ServeHTTP(rw, req)
}

//-------------------

func main() {

	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()

	// Install and instantiate the chaincode
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}

	app := Application{
		Fabric: &fSetup,
	}

	fmt.Println("finish chaincode declaration")

	fmt.Println("start server in port " + port)
	router := mux.NewRouter().StrictSlash(true)

	/* Authentication */
	router.HandleFunc("/auth", authentication).Methods("POST")
	router.HandleFunc("/createInitialUser", createInititalUser).Methods("POST")

	/* Users Routes */
	router.Handle("/users", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createUsersEndPoint)))).Methods("POST")
	router.Handle("/users", middleware.AuthMiddleware(http.HandlerFunc(allUsersEndPoint))).Methods("GET")
	router.Handle("/users/{id}", middleware.AuthMiddleware(http.HandlerFunc(findUserEndpoint))).Methods("GET")
	router.Handle("/users/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeUserEndpoint))).Methods("DELETE")
	router.Handle("/users/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateUserEndPoint)))).Methods("PUT")

	/* fileUpload */

	router.Handle("/fileUpload", middleware.AuthMiddleware(http.HandlerFunc(fileUpload))).Methods("POST")
	router.HandleFunc("/serveImage/{image}", serveImage).Methods("GET")
	router.Handle("/deleteFile/{file}", middleware.AuthMiddleware(http.HandlerFunc(deleteImage))).Methods("DELETE")
	router.Handle("/downloadFile/{file}", middleware.AuthMiddleware(http.HandlerFunc(downloadFile))).Methods("GET")

	/* testing chaincode */
	router.HandleFunc("/queryHelloChainCode", app.queryHelloChainCode).Methods("GET")
	router.HandleFunc("/invokeHelloChainCode/{word}", app.invokeHelloChaincode).Methods("GET")
	router.HandleFunc("/historyHelloChainCode", app.queryHelloChainCode).Methods("GET")

	/* Participants */
	//router.HandleFunc("/participants", authentication).Methods("GET")

	/* Assets */
	//router.HandleFunc("/assets", authentication).Methods("GET")

	/* ISSUES */
	//router.HandleFunc("/issues", authentication).Methods("GET")

	/* TRAZABILITY */
	//router.HandleFunc("/traz/{id}", authentication).Methods("GET")

	//start server
	log.Fatal(http.ListenAndServe(":"+port, &CORSRouterDecorator{router}))

}
