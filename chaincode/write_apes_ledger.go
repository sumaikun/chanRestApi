package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *ApesChainCode) saveParticipant(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesChainCode create Participant ###########")
	var err error

	// 0 ,    1,           2,       3,     4,             5,        6
	// name,  nationality, address, phone, identification, photoUrl, description

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	if len(args[0]) <= 0 {
		return shim.Error("name argument must be a non-empty string")
	}

	if len(args[1]) <= 0 {
		return shim.Error("nationality argument must be a non-empty string")
	}

	if len(args[2]) <= 0 {
		return shim.Error("address argument must be a non-empty string")
	}

	if len(args[3]) <= 0 {
		return shim.Error("phone argument must be a non-empty string")
	}

	if len(args[4]) <= 0 {
		return shim.Error("identification argument must be a non-empty string")
	}

	if len(args[6]) <= 0 {
		return shim.Error("description argument must be a non-empty string")
	}

	name := strings.ToLower(args[0])

	address := strings.ToLower(args[2])

	identification := strings.ToLower(args[4])

	nationality := strings.ToLower(args[1])

	objectType := "participant"

	phone := args[3]

	photoURL := args[5]

	description := args[6]

	participantAsBytes, err := stub.GetState(identification)
	if err != nil {
		return shim.Error("Failed to get participant: " + err.Error())
	} else if participantAsBytes != nil {
		fmt.Println("This participant already exists: " + identification)
		//return shim.Error("This participant already exists: " + identification)
	}

	participant := &Participant{objectType, name, nationality, address, phone, identification, photoURL, description}
	participantJSONasBytes, err := json.Marshal(participant)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(identification, participantJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("saveParticipant", []byte(identification))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}
