package main

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *ApesWallet) getData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesWallet get DATA by key ###########")

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

func (t *ApesWallet) getObjectTypeWithKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesWallet get Object Type results by key ###########")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting objectType to query")
	}

	docType := args[0]

	resultsIterator, err := stub.GetStateByPartialCompositeKey("type~identification", []string{docType})

	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	var i int

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false

	for i = 0; resultsIterator.HasNext(); i++ {

		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		//fmt.Println(responseRange)

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		if err != nil {
			return shim.Error(err.Error())
		}

		returnedType := compositeKeyParts[0]
		returnedID := compositeKeyParts[1]

		fmt.Println(objectType)

		fmt.Println(returnedType)

		fmt.Println(returnedID)

		valAsbytes, err := stub.GetState(returnedID)

		buffer.WriteString(string(valAsbytes))

		bArrayMemberAlreadyWritten = true

	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())

}

func (t *ApesWallet) getObjectTypeByKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesWallet get Object Type results by key ###########")

	// 0 , 1
	// key, objectType

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting key and objectType to query")
	}

	key := args[1]

	docType := args[1]

	resultsIterator, err := stub.GetStateByPartialCompositeKey("type~identification", []string{docType})

	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	var i int

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false

	for i = 0; resultsIterator.HasNext(); i++ {

		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		//fmt.Println(responseRange)

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		if err != nil {
			return shim.Error(err.Error())
		}

		returnedType := compositeKeyParts[0]
		returnedID := compositeKeyParts[1]

		fmt.Println(objectType)

		fmt.Println(returnedType)

		fmt.Println(returnedID)

		valAsbytes, err := stub.GetState(returnedID)

		buffer.WriteString(string(valAsbytes))

		bArrayMemberAlreadyWritten = true

	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())

}
