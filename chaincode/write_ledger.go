package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// invoke
// Every functions that read and write in the ledger will be here
func (t *ApesChainCode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesChainCode invoke ###########")

	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Check if the ledger key is "hello" and process if it is the case. Otherwise it returns an error.
	if args[0] == "hello" && len(args) == 2 {

		// Write the new value in the ledger
		err := stub.PutState("hello", []byte(args[1]))
		if err != nil {
			return shim.Error("Failed to update state of hello")
		}

		// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
		err = stub.SetEvent("eventInvoke", []byte{})
		if err != nil {
			return shim.Error(err.Error())
		}

		// Return this value in response
		return shim.Success(nil)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown invoke action, check the second argument.")
}
