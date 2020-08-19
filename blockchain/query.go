package blockchain

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// QueryHello query the chaincode to get the state of hello
func (setup *FabricSetup) QueryHello() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "query")
	args = append(args, "hello")

	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		fmt.Errorf("failed to query: %v", err)
		return "", err
	}

	return string(response.Payload), nil
}

// QueryGetData query the chaincode to get the state of a key
func (setup *FabricSetup) QueryGetData(key string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "getData")
	args = append(args, key)

	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}})
	if err != nil {
		fmt.Errorf("failed to query: %v", err)
		return "", err
	}

	return string(response.Payload), nil
}
