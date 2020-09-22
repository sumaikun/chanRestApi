package blockchain

import (
	"fmt"
	//"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// HistoryHello method example
func (setup *FabricSetup) HistoryHello() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "history")
	args = append(args, "hello")

	fmt.Println("args to send")
	fmt.Println(args)

	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return "", fmt.Errorf("failed to query history: %v", err)
	}

	return string(response.Payload), nil
}

// HistoryKey method example
func (setup *FabricSetup) HistoryKey(key string) ([]byte, error) {

	// Prepare arguments
	var args []string
	args = append(args, "historyKey")
	args = append(args, key)

	fmt.Println("args to send")
	fmt.Println(args)

	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %v", err)
	}

	return response.Payload, nil
}

//--------------------------------------------------- wallet driver -------------------------------------------------------------------//

// HistoryKey2 method example
func (setup *FabricSetup) HistoryKey2(key string) ([]byte, error) {

	// Prepare arguments
	var args []string
	args = append(args, "historyKey")
	args = append(args, key)

	fmt.Println("args to send")
	fmt.Println(args)

	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID2, Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %v", err)
	}

	return response.Payload, nil
}
