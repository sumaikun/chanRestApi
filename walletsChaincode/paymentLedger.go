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

	paymentType := string(args[4])

	identification := strings.ToLower(args[5])

	anumTypes := []string{"PAY", "DISCOUNT"}

	containtType := Contains(anumTypes, paymentType)

	if containtType != true {
		return shim.Error("The " + paymentType + " field must be a valid value for payment type Enum")
	}

	if paymentType == "PAY" {

		updateOwner.Balance = updateOwner.Balance + quantity

		objectType := "externalPayment"
		state := "success"

		externalPayment := &ExternalPayment{objectType, fromExternal, toWallet, state, date, quantity, paymentType, identification}
		externalPaymentJSONasBytes, err := json.Marshal(externalPayment)
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

		indexName = "wallet~externalPayment"
		typeIndexKey, err = stub.CreateCompositeKey(indexName, []string{updateOwner.Identification, identification})
		if err != nil {
			return shim.Error(err.Error())
		}
		value = []byte{0x00}
		stub.PutState(typeIndexKey, value)

		indexName = "externalAgent~externalPayment"
		typeIndexKey, err = stub.CreateCompositeKey(indexName, []string{externalAgentObject.Identification, identification})
		if err != nil {
			return shim.Error(err.Error())
		}
		value = []byte{0x00}
		stub.PutState(typeIndexKey, value)

	} else {
		if updateOwner.Balance < quantity {
			objectType := "externalPayment"

			state := "failed"

			externalPayment := &ExternalPayment{objectType, fromExternal, toWallet, state, date, quantity, paymentType, identification}
			externalPaymentJSONasBytes, err := json.Marshal(externalPayment)
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

			indexName = "wallet~externalPayment"
			typeIndexKey, err = stub.CreateCompositeKey(indexName, []string{updateOwner.Identification, identification})
			if err != nil {
				return shim.Error(err.Error())
			}
			value = []byte{0x00}
			stub.PutState(typeIndexKey, value)

			indexName = "externalAgent~externalPayment"
			typeIndexKey, err = stub.CreateCompositeKey(indexName, []string{externalAgentObject.Identification, identification})
			if err != nil {
				return shim.Error(err.Error())
			}
			value = []byte{0x00}
			stub.PutState(typeIndexKey, value)

			return shim.Error("not enought fonds  for this external payment")
		}

		updateOwner.Balance = updateOwner.Balance - quantity

		objectType := "externalPayment"

		state := "success"

		externalPayment := &ExternalPayment{objectType, fromExternal, toWallet, state, date, quantity, paymentType, identification}
		externalPaymentJSONasBytes, err := json.Marshal(externalPayment)
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

		indexName = "wallet~externalPayment"
		typeIndexKey, err = stub.CreateCompositeKey(indexName, []string{updateOwner.Identification, identification})
		if err != nil {
			return shim.Error(err.Error())
		}
		value = []byte{0x00}
		stub.PutState(typeIndexKey, value)

		indexName = "externalAgent~externalPayment"
		typeIndexKey, err = stub.CreateCompositeKey(indexName, []string{externalAgentObject.Identification, identification})
		if err != nil {
			return shim.Error(err.Error())
		}
		value = []byte{0x00}
		stub.PutState(typeIndexKey, value)

	}

	return shim.Success(nil)

}

func (t *ApesWallet) makeWalletPayment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### ApesWallet makeWalletPayment ###########")
	var err error

	// 0 ,             1,        2,      3,        4,
	// FromWallet,  ToWallet, date, quantity,  identification

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	if len(args[0]) <= 0 {
		return shim.Error("fromWallet argument must be a non-empty string")
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
		return shim.Error("identification argument must be a non-empty string")
	}

	fromWallet := args[0]

	fromWalletAsBytes, err := stub.GetState(fromWallet)
	if err != nil {
		return shim.Error("Failed to get  from wallet: " + err.Error())
	} else if fromWalletAsBytes == nil {
		return shim.Error("This wallet dont exists: " + fromWallet)
	}

	toWallet := args[1]

	toWalletAsBytes, err := stub.GetState(toWallet)
	if err != nil {
		return shim.Error("Failed to get reception owner wallet: " + err.Error())
	} else if toWalletAsBytes == nil {
		//fmt.Println("This externalAgent already exists: " + identification)
		return shim.Error("This owner wallet dont exists: " + toWallet)
	}

	updateFromOwner := Owner{}
	_ = json.Unmarshal(fromWalletAsBytes, &updateFromOwner)

	updateToOwner := Owner{}
	_ = json.Unmarshal(toWalletAsBytes, &updateToOwner)

	date := strings.ToLower(args[2])

	quantity, err := strconv.Atoi(args[3])

	if err != nil {
		return shim.Error("quantity argument must be a numeric string")
	}

	identification := strings.ToLower(args[4])

	if updateFromOwner.Balance < quantity {

		walletPayment := &WalletPayment{"walletPayment", fromWallet, toWallet, "failed", date, quantity, identification}
		walletPaymentJSONasBytes, err := json.Marshal(walletPayment)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = stub.PutState(identification, walletPaymentJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		indexName := "type~identification"
		typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{"walletPayment", identification})
		if err != nil {
			return shim.Error(err.Error())
		}
		value := []byte{0x00}
		stub.PutState(typeIndexKey, value)

		indexName = "fromWallet~walletPayment"
		typeIndexKey, err = stub.CreateCompositeKey(indexName, []string{updateFromOwner.Identification, identification})
		if err != nil {
			return shim.Error(err.Error())
		}
		value = []byte{0x00}
		stub.PutState(typeIndexKey, value)

		indexName = "toWallet~walletPayment"
		typeIndexKey, err = stub.CreateCompositeKey(indexName, []string{updateToOwner.Identification, identification})
		if err != nil {
			return shim.Error(err.Error())
		}
		value = []byte{0x00}
		stub.PutState(typeIndexKey, value)

		return shim.Error("not enought fonds  for this wallet payment")
	}

	updateFromOwner.Balance = updateFromOwner.Balance - quantity

	walletPayment := &WalletPayment{"walletPayment", fromWallet, toWallet, "failed", date, quantity, identification}
	walletPaymentJSONasBytes, err := json.Marshal(walletPayment)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(identification, walletPaymentJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	updateFromOwnerJSONasBytes, err := json.Marshal(updateFromOwner)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(updateFromOwner.Identification, updateFromOwnerJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	updateToOwner.Balance = updateToOwner.Balance + quantity

	updateToOwnerJSONasBytes, err := json.Marshal(updateToOwner)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(updateToOwner.Identification, updateToOwnerJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "type~identification"
	typeIndexKey, err := stub.CreateCompositeKey(indexName, []string{"walletPayment", identification})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(typeIndexKey, value)

	indexName = "fromWallet~walletPayment"
	typeIndexKey, err = stub.CreateCompositeKey(indexName, []string{updateFromOwner.Identification, identification})
	if err != nil {
		return shim.Error(err.Error())
	}
	value = []byte{0x00}
	stub.PutState(typeIndexKey, value)

	indexName = "toWallet~walletPayment"
	typeIndexKey, err = stub.CreateCompositeKey(indexName, []string{updateToOwner.Identification, identification})
	if err != nil {
		return shim.Error(err.Error())
	}
	value = []byte{0x00}
	stub.PutState(typeIndexKey, value)

	return shim.Success(nil)

}
