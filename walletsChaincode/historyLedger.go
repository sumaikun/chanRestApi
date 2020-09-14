package main

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *ApesWallet) historyKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesWallet history ###########")

	var key string
	// Check whether the number of arguments is sufficient
	if len(args) != 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	key = args[0]

	//historyData := []string{}

	historyIter, err := stub.GetHistoryForKey(key)

	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] cannot retrieve history for key <%s>, due to %s", key, err)
		fmt.Println(errMsg)
		return shim.Error(errMsg)
	}

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false

	for historyIter.HasNext() {
		modification, err := historyIter.Next()
		if err != nil {
			errMsg := fmt.Sprintf("[ERROR] cannot read record modification for key %s, id <%s>, due to %s", key, err)
			fmt.Println(errMsg)
			return shim.Error(errMsg)
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		fmt.Println("Returning information about", string(modification.Value))

		buffer.WriteString(string(modification.Value))
		//historyData = append(historyData, string(modification.Value))
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
