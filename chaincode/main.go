package main

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ApesChainCode implementation of Chaincode
type ApesChainCode struct {
}

// Owner representation in chaincode
type Owner struct {
	name           string
	nationality    string
	address        string
	phone          string
	identification string
	photoURL       string
	description    string
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

	// Put in the ledger the key/value hello/world
	err := stub.PutState("hello", []byte("world"))
	if err != nil {
		return shim.Error(err.Error())
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

	// Check whether it is an invoke request
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// In order to manage multiple type of request, we will check the first argument.
	// Here we have one possible argument: query (every query request will read in the ledger without modification)
	if args[0] == "query" {
		return t.query(stub, args)
	}

	// The update argument will manage all update in the ledger
	if args[0] == "invoke" {
		return t.invoke(stub, args)
	}

	// Check historian
	if args[0] == "history" {
		return t.history(stub, args)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument: " + args[0])
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(ApesChainCode))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
	}
}

//https://codeburst.io/writing-chaincode-in-golang-the-oop-way-4be3bb261dae
