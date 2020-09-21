package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	//"github.com/chainHero/heroes-service/blockchain"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sumaikun/apeslogistic-rest-api/blockchain"
	Config "github.com/sumaikun/apeslogistic-rest-api/config"
	Dao "github.com/sumaikun/apeslogistic-rest-api/dao"
	Helpers "github.com/sumaikun/apeslogistic-rest-api/helpers"
	middleware "github.com/sumaikun/apeslogistic-rest-api/middlewares"
	"github.com/sumaikun/apeslogistic-rest-api/models"
	"github.com/thedevsaddam/govalidator"
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var fSetup = blockchain.FabricSetup{
	// Network parameters
	OrdererID: "orderer.hf.chainhero.io",

	// Channel parameters
	ChannelID:     "chainhero",
	ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/chainHero/heroes-service/fixtures/artifacts/chainhero.channel.tx",

	// Chaincode parameters
	ChainCodeID:     "heroes-service",
	ChainCodeID2:    "heroes-wallet",
	ChaincodeGoPath: os.Getenv("GOPATH"),
	ChaincodePath:   "github.com/sumaikun/apeslogistic-rest-api/walletsChaincode/",
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

	govalidator.AddCustomRule("ifExist", func(field string, rule string, message string, value interface{}) error {

		if len(value.(string)) >= 0 {
			return nil
		}

		return fmt.Errorf("The %s field must exist", field)

	})

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
	router.HandleFunc("/historyHelloChainCode", app.historyHelloChainCode).Methods("GET")

	/*get data from chaincode */
	router.HandleFunc("/getChaincodeData/{key}", app.getDataFromChaincode).Methods("GET")
	router.Handle("/getHistoryForKey/{key}", middleware.AuthMiddleware(http.HandlerFunc(app.getHistoryForKey))).Methods("GET")

	/* Participants */
	router.Handle("/participants", middleware.AuthMiddleware(http.HandlerFunc(app.saveParticipant))).Methods("POST")
	router.Handle("/participants", middleware.AuthMiddleware(http.HandlerFunc(app.getParticipants))).Methods("GET")

	/* Assets */
	router.Handle("/assets", middleware.AuthMiddleware(http.HandlerFunc(app.saveAsset))).Methods("POST")
	router.Handle("/assets", middleware.AuthMiddleware(http.HandlerFunc(app.getAssets))).Methods("GET")

	/* Infrastructure */
	router.Handle("/installChainCode", middleware.AuthMiddleware(http.HandlerFunc(app.installChainCode))).Methods("POST")
	router.Handle("/instantiateChainCode", middleware.AuthMiddleware(http.HandlerFunc(app.instantiateChainCode))).Methods("GET")

	/*get data from  wallet chaincode */
	router.HandleFunc("/getChaincodeData2/{key}", app.getDataFromChaincode2).Methods("GET")
	router.Handle("/getHistoryForKey2/{key}", middleware.AuthMiddleware(http.HandlerFunc(app.getHistoryForKey2))).Methods("GET")

	/* Owners */
	router.Handle("/walletOwners", middleware.AuthMiddleware(middleware.UserMiddleware(middleware.OnlyAdminMiddleware(http.HandlerFunc(app.saveOwner))))).Methods("POST")
	router.Handle("/walletOwners", middleware.CognitoMiddleware(middleware.AuthMiddleware(http.HandlerFunc(app.getOwners)))).Methods("GET")

	/* External Agents */
	router.Handle("/walletExternalAgents", middleware.AuthMiddleware(middleware.UserMiddleware(middleware.OnlyAdminMiddleware(http.HandlerFunc(app.saveExternalAgent))))).Methods("POST")
	router.Handle("/walletExternalAgents", middleware.CognitoMiddleware(middleware.AuthMiddleware(app.CreateWalletIfNotExist(http.HandlerFunc(app.getExternalAgents))))).Methods("GET")

	/* Wallets Events */
	router.Handle("/walletEvents", middleware.AuthMiddleware(middleware.UserMiddleware(middleware.OnlyAdminMiddleware(http.HandlerFunc(app.saveEvent))))).Methods("POST")
	router.Handle("/walletEvents", middleware.AuthMiddleware(http.HandlerFunc(app.getEvents))).Methods("GET")

	/* Wallets Rules */
	router.Handle("/walletRules", middleware.AuthMiddleware(middleware.UserMiddleware(middleware.OnlyAdminMiddleware(http.HandlerFunc(app.saveRule))))).Methods("POST")
	router.Handle("/walletRules", middleware.AuthMiddleware(http.HandlerFunc(app.getRules))).Methods("GET")

	/* Wallet Payments */
	router.Handle("/externalPayment", middleware.AuthMiddleware(http.HandlerFunc(app.externalPayment))).Methods("POST")
	router.Handle("/walletPayment", middleware.AuthMiddleware(http.HandlerFunc(app.walletPayment))).Methods("POST")

	/* Trazability */
	router.Handle("/walletExternalPayment/{key}", middleware.AuthMiddleware(http.HandlerFunc(app.walletExternalPayment))).Methods("GET")
	router.Handle("/externalAgentExternalPayment/{key}", middleware.AuthMiddleware(http.HandlerFunc(app.walletExternalPayment))).Methods("GET")
	router.Handle("/fromWalletWalletPayment/{key}", middleware.AuthMiddleware(http.HandlerFunc(app.fromWalletWalletPayment))).Methods("GET")
	router.Handle("/toWalletWalletPayment/{key}", middleware.AuthMiddleware(http.HandlerFunc(app.toWalletWalletPayment))).Methods("GET")

	/* ISSUES */
	//router.HandleFunc("/issues", authentication).Methods("GET")

	/* TRAZABILITY */
	//router.HandleFunc("/traz/{id}", authentication).Methods("GET")

	//start server
	log.Fatal(http.ListenAndServe(":"+port, &CORSRouterDecorator{router}))

}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// CreateWalletIfNotExist Verify Wallet Existence
func (app *Application) CreateWalletIfNotExist(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cognitoEmail := context.Get(r, "cognito_email")

		if cognitoEmail != nil {
			cognitoEmailParsed := cognitoEmail.(*string)

			fmt.Println("cognitoEmailParsed", *cognitoEmailParsed)

			response, err := app.Fabric.QueryGetData2(*cognitoEmailParsed)
			if err != nil {
				fmt.Printf("Unable to query  the chaincode: %v\n", err)

				if strings.Contains(err.Error(), "Key does not exist") {
					fmt.Println("key does not exist then create")

					txID, err2 := app.Fabric.SaveOwner(models.Owner{"owner", *cognitoEmailParsed, "", "", "", *cognitoEmailParsed, *cognitoEmailParsed, "", "", 1})
					if err2 != nil {
						fmt.Printf("Unable to save owner on the chaincode: %v\n", err2)
						Helpers.RespondWithJSON(w, http.StatusBadGateway, map[string]string{"error": err2.Error()})
						return
					}

					next.ServeHTTP(w, r)

					return
				}

				Helpers.RespondWithJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
				return
			}

			fmt.Println("Response from chaincode: %s\n", response)

			next.ServeHTTP(w, r)

			return
		}

		Helpers.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "wallet not exist"})
		return

	})
}
