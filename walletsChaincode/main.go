package main

import (
	"fmt"
	"math/rand"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// ApesWallet implementation of Chaincode
type ApesWallet struct {
}

// Owner representation in chaincode
type Owner struct {
	ObjectType     string `json:"docType"`
	Name           string `json:"name"`
	Nationality    string `json:"nationality"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Identification string `json:"identification"`
	PhotoURL       string `json:"photoUrl"`
	Notes          string `json:"notes"`
	Balance        int    `json:"balance"`
}

// ExternalAgent representation in chaincode
type ExternalAgent struct {
	ObjectType     string `json:"docType"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Identification string `json:"identification"`
}

// WalletPayment representation in chaincode
type WalletPayment struct {
	ObjectType     string `json:"docType"`
	FromWallet     string `json:"fromWallet"`
	ToWallet       string `json:"toWallet"`
	State          string `json:"state"`
	Date           string `json:"date"`
	Quantity       int    `json:"quantity"`
	Identification string `json:"identification"`
}

// ExternalPayment representation in chaincode
type ExternalPayment struct {
	ObjectType     string `json:"docType"`
	FromExternal   string `json:"fromExternal"`
	ToWallet       string `json:"toWallet"`
	State          string `json:"state"`
	Date           string `json:"date"`
	Quantity       int    `json:"quantity"`
	PaymentType    string `json:"paymentType"`
	Identification string `json:"identification"`
}

// Event representation in chaincode
type Event struct {
	ObjectType     string `json:"docType"`
	FromExternal   string `json:"fromExternal"`
	FromWallet     string `json:"fromWallet"`
	ToWallet       string `json:"toWallet"`
	ToExternal     string `json:"toExternal"`
	Identification string `json:"identification"`
}

// Rule representation in chaincode
type Rule struct {
	ObjectType     string `json:"docType"`
	Event          string `json:"event"`
	Fee            int    `json:"fee"`
	ToWallet       string `json:"toWallet"`
	ToExternal     string `json:"toExternal"`
	Date           string `json:"date"`
	Quantity       int    `json:"quantity"`
	State          bool   `json:"state"`
	Identification string `json:"identification"`
}

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *ApesWallet) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### ApesWallet Init with changes ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call")
	}

	// Return a successful message
	return shim.Success(nil)
}

// Invoke
// All future requests named invoke will arrive here.
func (t *ApesWallet) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### ApesWallet Invoke ###########")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	fmt.Println("function:" + function)

	fmt.Println("args")

	fmt.Println(args)

	// Check historian
	if function == "historyKey" {
		return t.historyKey(stub, args)
	}

	// Check key
	if function == "getData" {
		return t.getData(stub, args)
	}

	// getObjectTypeWithKey
	if function == "getObjectTypeWithKey" {
		return t.getObjectTypeWithKey(stub, args)
	}

	// getObjectTypeByKey
	if function == "getObjectTypeByKey" {
		return t.getObjectTypeByKey(stub, args)
	}

	// saveOwner
	if function == "saveOwner" {
		return t.saveOwner(stub, args)
	}

	// saveExternalAgent
	if function == "saveExternalAgent" {
		return t.saveExternalAgent(stub, args)
	}

	// saveEvent
	if function == "saveEvent" {
		return t.saveEvent(stub, args)
	}

	// saveRule
	if function == "saveRule" {
		return t.saveRule(stub, args)
	}

	// makeExternalPayment
	if function == "makeExternalPayment" {
		return t.makeExternalPayment(stub, args)
	}

	// makeWalletPayment
	if function == "makeWalletPayment" {
		return t.makeWalletPayment(stub, args)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument: " + args[0])
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(ApesWallet))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
	}
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// RandStringRunes for generate random string
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
