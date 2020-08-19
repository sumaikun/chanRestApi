package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *ApesChainCode) getData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesChainCode get Participant ###########")

	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key) //get the key from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Key does not exist: " + key + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)

}
