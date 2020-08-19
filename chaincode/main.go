package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ApesChainCode implementation of Chaincode
type ApesChainCode struct {
}

// Participant representation in chaincode
type Participant struct {
	ObjectType     string `json:"docType"`
	Name           string `json:"name"`
	Nationality    string `json:"nationality"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	Identification string `json:"identification"`
	PhotoURL       string `json:"photoUrl"`
	Description    string `json:"description"`
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
	/*if function == "invoke" {
		return shim.Error("Unknown function call")
	}*/

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	fmt.Println("function:" + function)

	fmt.Println("args")

	fmt.Println(args)

	// In order to manage multiple type of request, we will check the first argument.
	// Here we have one possible argument: query (every query request will read in the ledger without modification)
	if function == "query" {
		return t.query(stub, args)
	}

	// The update argument will manage all update in the ledger
	if function == "invoke" {
		return t.invoke(stub, args)
	}

	// Check historian
	if function == "history" {
		return t.history(stub, args)
	}

	// Check key
	if function == "getData" {
		return t.getData(stub, args)
	}

	// Create participant
	if function == "createParticipant" {
		return t.createParticipant(stub, args)
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
