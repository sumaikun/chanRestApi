package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *ApesWallet) saveOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesWallet create Owner ###########")
	var err error

	// 0 ,    1,           2,       3,     4,             5,        6      7
	// name,  nationality, address, phone, identification, photoUrl, notes, balance

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
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

	name := strings.ToLower(args[0])

	address := strings.ToLower(args[2])

	identification := strings.ToLower(args[4])

	nationality := strings.ToLower(args[1])

	objectType := "owner"

	phone := args[3]

	photoURL := args[5]

	notes := args[6]

	balance := 0

	ownerAsBytes, err := stub.GetState(identification)
	if err != nil {
		return shim.Error("Failed to get owner: " + err.Error())
	} else if ownerAsBytes != nil {
		fmt.Println("This owner already exists: " + identification)
		//return shim.Error("This owner already exists: " + identification)
	}

	owner := &Owner{objectType, name, nationality, address, phone, identification, photoURL, notes}
	ownerJSONasBytes, err := json.Marshal(owner)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(identification, ownerJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("saveOwner", []byte(identification))
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "type~identification"
	typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{objectType, identification})
	if err != nil {
		return shim.Error(err.Error())
	}

	value := []byte{0x00}
	stub.PutState(typeIndexKey, value)

	return shim.Success(nil)

}

func (t *ApesWallet) saveExternalAgent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesWallet create ExternalAgent ###########")
	var err error

	// 0 ,    1,          2
	// name,  description, identification

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	if len(args[0]) <= 0 {
		return shim.Error("name argument must be a non-empty string")
	}

	if len(args[1]) <= 0 {
		return shim.Error("description argument must be a non-empty string")
	}

	if len(args[2]) <= 0 {
		return shim.Error("identification argument must be a non-empty string")
	}

	name := strings.ToLower(args[0])

	description := strings.ToLower(args[1])

	identification := strings.ToLower(args[2])

	objectType := "externalAgent"

	externalAgentAsBytes, err := stub.GetState(identification)
	if err != nil {
		return shim.Error("Failed to get externalAgent: " + err.Error())
	} else if externalAgentAsBytes != nil {
		fmt.Println("This externalAgent already exists: " + identification)
		//return shim.Error("This externalAgent already exists: " + identification)
	}

	externalAgent := &ExternalAgent{objectType, name, description, identification}
	externalAgentJSONasBytes, err := json.Marshal(externalAgent)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(identification, externalAgentJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("saveExternalAgent", []byte(identification))
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "type~identification"
	typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{objectType, identification})
	if err != nil {
		return shim.Error(err.Error())
	}

	value := []byte{0x00}
	stub.PutState(typeIndexKey, value)

	return shim.Success(nil)

}

func (t *ApesWallet) saveEvent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 0,            1,           2,       3,
	// fromExternal, fromWallet, toWallet, toExternal

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	if len(args[0]) > 0 && len(args[1]) > 0 {
		return shim.Error("its only valid from an external argent or a wallet not both for input event")
	}

	if len(args[0]) <= 0 && len(args[1]) <= 0 {
		return shim.Error("must give  an external argent or a wallet as input event")
	}

	if len(args[2]) <= 0 && len(args[3]) <= 0 {
		return shim.Error("must give  an external argent or a wallet as output event")
	}

	if len(args[2]) > 0 && len(args[3]) > 0 {
		return shim.Error("its only valid  an external argent or a wallet not both for output event")
	}

	fromExternal := strings.ToLower(args[0])

	fromWallet := strings.ToLower(args[1])

	toWallet := strings.ToLower(args[2])

	toExternal := strings.ToLower(args[3])

	objectType := "event"

	var keyEvent string

	if len(fromExternal) > 0 && len(toWallet) > 0 {
		keyEvent = "event-" + fromExternal + "-" + toWallet
	}

	if len(fromExternal) > 0 && len(toExternal) > 0 {
		keyEvent = "event-" + fromExternal + "-" + toExternal
	}

	if len(fromWallet) > 0 && len(toExternal) > 0 {
		keyEvent = "event-" + fromWallet + "-" + toExternal
	}

	if len(fromWallet) > 0 && len(toWallet) > 0 {
		keyEvent = "event-" + fromWallet + "-" + toWallet
	}

	eventAsBytes, err := stub.GetState(keyEvent)

	if err != nil {
		return shim.Error("Failed to get event: " + err.Error())
	} else if eventAsBytes != nil {
		fmt.Println("This event already exists: " + keyEvent)
		//return shim.Error("This externalAgent already exists: " + identification)
	}

	event := &Event{objectType, name, description, identification}
	eventJSONasBytes, err := json.Marshal(event)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(keyEvent, eventJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("saveEvent", []byte(keyEvent))
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "type~identification"
	typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{objectType, keyEvent})
	if err != nil {
		return shim.Error(err.Error())
	}

	value := []byte{0x00}
	stub.PutState(typeIndexKey, value)

	return shim.Success(nil)

}

func (t *ApesWallet) saveRule(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 0,        1,     2,       3,       4,    5,       6
	// event, fee, toWallet, toExternal, date, quantity, state

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	if len(args[2]) > 0 && len(args[3]) > 0 {
		return shim.Error("its only valid  an external argent or a wallet not both for output rule")
	}

	if len(args[2]) == 0 && len(args[3]) == 0 {
		return shim.Error("must give  an external argent or a wallet for an output rule")
	}

	if len(args[1]) > 0 && len(args[5]) > 0 {
		return shim.Error("its only valid  a fee or a quantity for output rule")
	}

	if len(args[1]) == 0 && len(args[5]) == 0 {
		return shim.Error("Must suminister a fee or a quantity")
	}

	if len(args[0]) <= 0 {
		return shim.Error("event argument must be a non-empty string")
	}

	if len(args[4]) <= 0 {
		return shim.Error("date argument must be a non-empty string")
	}

	if len(args[6]) <= 0 {
		return shim.Error("state argument must be a non-empty string")
	}

	event := strings.ToLower(args[0])

	fee, err := strconv.Atoi(args[1])

	if err != nil {
		return shim.Error("quantity argument must be a numeric string")
	}

	toWallet := strings.ToLower(args[2])

	toExternal := strings.ToLower(args[3])

	date := strings.ToLower(args[4])

	quantity, err := strconv.Atoi(args[5])

	if err != nil {
		return shim.Error("quantity argument must be a numeric string")
	}

	state, err := strconv.ParseBool(args[6])
	if err != nil {
		return shim.Error("6th argument must be a boolean string")
	}

	objectType := "rule"

	var keyRule string

	if len(toWallet) > 0 {
		keyRule = event + "-rule-" + toWallet
	}

	if len(toExternal) > 0 {
		keyRule = event + "-rule-" + toExternal
	}

	ruleAsBytes, err := stub.GetState(keyEvent)

	if err != nil {
		return shim.Error("Failed to get rule: " + err.Error())
	} else if ruleAsBytes != nil {
		fmt.Println("This rule already exists: " + keyRule)
		//return shim.Error("This externalAgent already exists: " + identification)
	}

	rule := &Rule{objectType, event, fee, toWallet, toExternal, date, quantity, state}

	ruleJSONasBytes, err := json.Marshal(rule)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(keyRule, ruleJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("saveRule", []byte(keyRule))
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "type~identification"
	typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{objectType, keyRule})
	if err != nil {
		return shim.Error(err.Error())
	}

	value := []byte{0x00}
	stub.PutState(typeIndexKey, value)

	indexName := "event~rule"
	typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{event, keyRule})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(typeIndexKey, value)

	return shim.Success(nil)

}
