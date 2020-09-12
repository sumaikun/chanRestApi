package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


func (t *ApesWallet) makeExternalPayment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesWallet makeExternalPayment ###########")
	var err error

	// 0 ,             1,          2  3        4     5
	// fromExternal,  toWallet, date, quantity, paymentType, identification

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	if len(args[0]) <= 0 {
		return shim.Error("fromExternal argument must be a non-empty string")
	}

	if len(args[1]) <= 0 {
		return shim.Error("toWallet argument must be a non-empty string")
	}

	if len(args[2]) <= 0 {
		return shim.Error("date argument must be a non-empty string")
	}

	if len(args[3]) <= 0 {
		return shim.Error("quantity argument must be a non-empty string")
	}

	if len(args[4]) <= 0 {
		return shim.Error("type argument must be a non-empty string")
	}

	if len(args[5]) <= 0 {
		return shim.Error("identification argument must be a non-empty string")
	}

	fromExternal := args[0]

	externalAgentAsBytes, err := stub.GetState(fromExternal)
	if err != nil {
		return shim.Error("Failed to get externalAgent: " + err.Error())
	} else if externalAgentAsBytes == nil {
		//fmt.Println("This externalAgent already exists: " + identification)
		return shim.Error("This externalAgent dont exists: " + fromExternal)
	}

	toWallet := args[1]

	ownerAsBytes, err := stub.GetState(toWallet)
	if err != nil {
		return shim.Error("Failed to get owner wallet: " + err.Error())
	} else if ownerAsBytes == nil {
		//fmt.Println("This externalAgent already exists: " + identification)
		return shim.Error("This owner wallet dont exists: " + toWallet)
	}

	externalAgentObject := ExternalAgent{}
	_ = json.Unmarshal(externalAgentAsBytes, &externalAgentObject)


	updateOwner := Owner{}
	_ = json.Unmarshal(ownerAsBytes, &updateOwner)

	date := args[2]

	quantity, err := strconv.Atoi(args[3])

	if err != nil {
		return shim.Error("quantity argument must be a numeric string")
	}

	paymentType := args[4]

	identification := args[5]

	anumTypes := []string{"PAY", "DISCOUNT"}

	containtType := Contains(paymentType)

	if containtType != true {
		return shim.Error("The " + containtType + " field must be a valid value for payment type Enum")
	}

	if containtType == "PAY"
	{
		updateOwner.Balance = updateOwner.Balance + quantity 

		externalPayment := &ExternalPayment{ "externalPayment", fromExternal, toWallet, "success", date, quantity, type, identification }
		externalPaymentJSONasBytes, err := json.Marshal(externalAgent)
		if err != nil {
			return shim.Error(err.Error())
		}		

		err = stub.PutState(identification, externalPaymentJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		updateOwnerJSONasBytes, err := json.Marshal(updateOwner)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = stub.PutState(updateOwner.Identification, updateOwnerJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		indexName := "type~identification"
		typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{"externalPayment", identification})
		if err != nil {
			return shim.Error(err.Error())
		}
		value := []byte{0x00}
		stub.PutState(typeIndexKey, value)

		indexName := "wallet~externalPayment"
		typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{updateOwner.Identification, identification})
		if err != nil {
			return shim.Error(err.Error())
		}
		value := []byte{0x00}
		stub.PutState(typeIndexKey, value)


		indexName := "externalAgent~externalPayment"
		typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{externalAgentObject.Identification, identification})
		if err != nil {
			return shim.Error(err.Error())
		}
		value := []byte{0x00}
		stub.PutState(typeIndexKey, value)

	}else{
		if updateOwner.Balance < quantity
		{
			externalPayment := &ExternalPayment{ "externalPayment", fromExternal, toWallet, "failed", date, quantity, type, identification }
			externalPaymentJSONasBytes, err := json.Marshal(externalAgent)
			if err != nil {
				return shim.Error(err.Error())
			}		

			err = stub.PutState(identification, externalPaymentJSONasBytes)
			if err != nil {
				return shim.Error(err.Error())
			}

			indexName := "type~identification"
			typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{"externalPayment", identification})
			if err != nil {
				return shim.Error(err.Error())
			}
			value := []byte{0x00}
			stub.PutState(typeIndexKey, value)

			indexName := "wallet~externalPayment"
			typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{updateOwner.Identification, identification})
			if err != nil {
				return shim.Error(err.Error())
			}
			value := []byte{0x00}
			stub.PutState(typeIndexKey, value)


			indexName := "externalAgent~externalPayment"
			typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{externalAgentObject.Identification, identification})
			if err != nil {
				return shim.Error(err.Error())
			}
			value := []byte{0x00}
			stub.PutState(typeIndexKey, value)
			
			return shim.Error("not enought fonds  for this external payment")
		}else{
			
			updateOwner.Balance = updateOwner.Balance - quantity 

			externalPayment := &ExternalPayment{ "externalPayment", fromExternal, toWallet, "success", date, quantity, type, identification }
			externalPaymentJSONasBytes, err := json.Marshal(externalAgent)
			if err != nil {
				return shim.Error(err.Error())
			}		

			err = stub.PutState(identification, externalPaymentJSONasBytes)
			if err != nil {
				return shim.Error(err.Error())
			}

			updateOwnerJSONasBytes, err := json.Marshal(updateOwner)
			if err != nil {
				return shim.Error(err.Error())
			}

			err = stub.PutState(updateOwner.Identification, updateOwnerJSONasBytes)
			if err != nil {
				return shim.Error(err.Error())
			}

			indexName := "type~identification"
			typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{"externalPayment", identification})
			if err != nil {
				return shim.Error(err.Error())
			}
			value := []byte{0x00}
			stub.PutState(typeIndexKey, value)

			indexName := "wallet~externalPayment"
			typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{updateOwner.Identification, identification})
			if err != nil {
				return shim.Error(err.Error())
			}
			value := []byte{0x00}
			stub.PutState(typeIndexKey, value)


			indexName := "externalAgent~externalPayment"
			typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{externalAgentObject.Identification, identification})
			if err != nil {
				return shim.Error(err.Error())
			}
			value := []byte{0x00}
			stub.PutState(typeIndexKey, value)

		}

	}

	return shim.Success(nil)

}