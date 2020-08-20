package main

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *ApesChainCode) getData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesChainCode get DATA by key ###########")

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

func (t *ApesChainCode) getObjectTypeWithKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesChainCode get Object Type results by key ###########")

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

	for i = 0; resultsIterator.HasNext(); i++ {

		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
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

	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())

}

/********************************** only works with couch db ***********************************/

func (t *ApesChainCode) getObjectType(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesChainCode get Object Type results ###########")

	var objectType string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting objectType to query")
	}

	objectType = args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\"}}", objectType)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)

}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

/********************************** only works with couch db ***********************************/
