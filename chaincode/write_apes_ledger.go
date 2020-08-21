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

	indexName := "type~identification"
	typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{objectType, identification})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(typeIndexKey, value)

	return shim.Success(nil)

}

func (t *ApesChainCode) saveAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesChainCode create Asset ###########")

	// 0 ,            1,           2,    3,     4,             5,    6    7
	// participant,  state, location, meta, identification, title, date, assetType

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}

	if len(args[0]) <= 0 {
		return shim.Error("participant argument must be a non-empty string")
	}

	if len(args[1]) <= 0 {
		return shim.Error("state argument must be a non-empty string")
	}

	if len(args[2]) <= 0 {
		return shim.Error("location argument must be a non-empty string")
	}

	if len(args[4]) <= 0 {
		return shim.Error("identification argument must be a non-empty string")
	}

	if len(args[5]) <= 0 {
		return shim.Error("title argument must be a non-empty string")
	}

	if len(args[6]) <= 0 {
		return shim.Error("date argument must be a non-empty string")
	}

	if len(args[7]) <= 0 {
		return shim.Error("assetType argument must be a non-empty string")
	}

	participant := strings.ToLower(args[0])

	participantAsBytes, err := stub.GetState(participant)
	if err != nil {
		return shim.Error("Failed to get participant: " + err.Error())
	} else if participantAsBytes == nil {
		//fmt.Println("This participant already exists: " + identification)
		return shim.Error("This participant does not exists: " + participant)
	}

	state := strings.ToLower(args[1])

	location := strings.ToLower(args[2])

	meta := strings.ToLower(args[3])

	identification := strings.ToLower(args[4])

	title := strings.ToLower(args[5])

	date := strings.ToLower(args[6])

	assetType := strings.ToLower(args[7])

	objectType := "asset"

	asset := &Asset{objectType, participant, state, location, meta, identification, title, date, assetType}
	assetJSONasBytes, err := json.Marshal(asset)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(identification, participantJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("saveAsset", []byte(identification))
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "type~identification"
	typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{objectType, identification})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(typeIndexKey, value)

	return shim.Success(nil)

}
