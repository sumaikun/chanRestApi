package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

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
	Identification string `json:"identification"`
	PhotoURL       string `json:"photoUrl"`
	Notes          string `json:"notes"`
	Balance        string `json:"balance"`
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
	FromWallet     string `json:"from"`
	ToWallet       string `json:"to"`
	State          string `json:"state"`
	Date           string `json:"date"`
	Quantity       string `json:"quantity"`
	Identification string `json:"identification"`
}

// ExternalPayment representation in chaincode
type ExternalPayment struct {
	ObjectType     string `json:"docType"`
	FromExternal   string `json:"from"`
	ToWallet       string `json:"to"`
	State          string `json:"state"`
	Date           string `json:"date"`
	Quantity       int    `json:"quantity"`
	PaymentType    string `json:"paymentType"`
	Identification string `json:"identification"`
}

// Event representation in chaincode
type Event struct {
	ObjectType   string `json:"docType"`
	fromExternal string `json:"fromExternal"`
	fromWallet   string `json:"fromWallet"`
	toWallet     string `json:"toWallet"`
	toExternal   string `json:"toExternal"`
}

// Rule representation in chaincode
type Rule struct {
	ObjectType string `json:"docType"`
	Event      string `json:"event"`
	fee        int    `json:"fee"`
	ToWallet   string `json:"toWallet"`
	toExternal string `json:"toAgent"`
	Date       string `json:"date"`
	Quantity   int    `json:"quantity"`
	State      bool   `json:"state"`
}

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *ApesChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### ApesChainCode Init with changes ###########")

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
func (t *ApesChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### ApesChainCode Invoke ###########")

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

	// Get data result from objectType
	if function == "getObjectType" {
		return t.getObjectType(stub, args)
	}

	// saveOwner
	if function == "getObjectType" {
		return t.saveOwner(stub, args)
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
		return t.saveEven(stub, args)
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
