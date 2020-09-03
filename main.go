package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	Config "github.com/sumaikun/apeslogistic-rest-api/config"
	Dao "github.com/sumaikun/apeslogistic-rest-api/dao"
	middleware "github.com/sumaikun/apeslogistic-rest-api/middlewares"
)

//Application object to chaincode connection
type Application struct {
	gway *gateway.Gateway
}

var (
	port   string
	jwtKey []byte
)

var dao = Dao.MongoConnector{}

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

	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			os.Exit(1)
		}
	}

	ccpPath := filepath.Join(
		"/",
		"home",
		"ubuntu",
		"go",
		"src",
		"github.com",
		"sumaikun",
		"apeslogistic-rest-api",
		"fixtures",
		"connection-Budget.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("airlinechannel")
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(network)

	contract := network.GetContract("gocc1")

	result, err := contract.EvaluateTransaction("query", "hello")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	app := Application{
		gway: gw,
	}

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
	router.HandleFunc("/historyHelloChainCode", app.historyHelloChainCode).Methods("GET")

	/*get data from chaincode */
	//router.Handle("/getChaincodeData/{key}", middleware.AuthMiddleware(http.HandlerFunc(app.getDataFromChaincode))).Methods("GET")

	/* Participants */
	//router.Handle("/participants", middleware.AuthMiddleware(http.HandlerFunc(app.saveParticipant))).Methods("POST")
	//router.Handle("/participants", middleware.AuthMiddleware(http.HandlerFunc(app.getParticipants))).Methods("GET")

	/* Assets */
	//router.Handle("/assets", middleware.AuthMiddleware(http.HandlerFunc(app.saveAsset))).Methods("POST")
	//router.Handle("/assets", middleware.AuthMiddleware(http.HandlerFunc(app.getAssets))).Methods("GET")

	/* Infrastructure */
	//router.Handle("/installChainCode", middleware.AuthMiddleware(http.HandlerFunc(app.installChainCode))).Methods("POST")
	//router.Handle("/instantiateChainCode", middleware.AuthMiddleware(http.HandlerFunc(app.instantiateChainCode))).Methods("GET")

	/* ISSUES */
	//router.HandleFunc("/issues", authentication).Methods("GET")

	/* TRAZABILITY */
	//router.HandleFunc("/traz/{id}", authentication).Methods("GET")

	//start server
	log.Fatal(http.ListenAndServe(":"+port, &CORSRouterDecorator{router}))

}

func populateWallet(wallet *gateway.Wallet) error {

	credPath := filepath.Join(
		"/",
		"home",
		"ubuntu",
		"go",
		"src",
		"github.com",
		"sumaikun",
		"apeslogistic-rest-api",
		"fixtures",
		"crypto-config",
		"peerOrganizations",
		"budget.com",
		"users",
		"User1@budget.com",
		"msp",
	)

	//credPath += "/" + credPath

	//fmt.Println(credPath)

	certPath := filepath.Join(credPath, "signcerts", "User1@budget.com-cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key

	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("BudgetMSP", string(cert), string(key))

	err = wallet.Put("appUser", identity)
	if err != nil {
		return err
	}
	return nil
}
