package main

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// query
// Every readonly functions in the ledger will be here
func (t *ApesChainCode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesChainCode query with cert ###########")

	creatorByte, err := stub.GetCreator()

	if err != nil {
		fmt.Println("failed to get Creator: %v", err)
	}

	fmt.Println("Creator ----------------------------------")

	//fmt.Println(creatorByte)

	fmt.Println(string(creatorByte))

	//fmt.Println("Proposal ----------------------------------")

	//fmt.Println(stub.GetSignedProposal())

	//fmt.Println("Transient ----------------------------------")

	//fmt.Println(stub.GetTransient())

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Like the Invoke function, we manage multiple type of query requests with the second argument.
	// We also have only one possible argument: hello
	if args[1] == "hello" {

		// Get the state of the value matching the key hello in the ledger
		state, err := stub.GetState("hello")
		if err != nil {
			return shim.Error("Failed to get state of hello")
		}

		// Return this value in response
		return shim.Success(state)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown query action, check the  argument.")
}

func (t *ApesChainCode) history(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesChainCode history ###########")

	// Check whether the number of arguments is sufficient
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Like the Invoke function, we manage multiple type of query requests with the second argument.
	// We also have only one possible argument: hello
	if args[1] == "hello" {

		key := "hello"

		historyData := []string{}

		historyIter, err := stub.GetHistoryForKey(key)

		if err != nil {
			errMsg := fmt.Sprintf("[ERROR] cannot retrieve history for key <%s>, due to %s", key, err)
			fmt.Println(errMsg)
			return shim.Error(errMsg)
		}

		for historyIter.HasNext() {
			modification, err := historyIter.Next()
			if err != nil {
				errMsg := fmt.Sprintf("[ERROR] cannot read record modification for key %s, id <%s>, due to %s", key, err)
				fmt.Println(errMsg)
				return shim.Error(errMsg)
			}
			fmt.Println("Returning information about", string(modification.Value))
			historyData = append(historyData, string(modification.Value))
		}

		historyAsString := strings.Join(historyData, ",")

		return shim.Success([]byte("[" + historyAsString + "]"))
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown history action, check the second argument.")
}
