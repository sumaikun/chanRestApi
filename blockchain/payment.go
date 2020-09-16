package blockchain

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/pkg/errors"

	Models "github.com/sumaikun/apeslogistic-rest-api/models"
)

// ExternalPayment driver
func (setup *FabricSetup) ExternalPayment(externalPayment Models.ExternalPayment) (string, error) {
	// Prepare arguments
	var args []string
	args = append(args, "makeExternalPayment")
	args = append(args, externalPayment.FromExternal)
	args = append(args, externalPayment.ToWallet)
	args = append(args, externalPayment.Date)
	args = append(args, strconv.Itoa(externalPayment.Quantity))
	args = append(args, externalPayment.PaymentType)
	args = append(args, externalPayment.Identification)

	eventID := "makeExternalPayment"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	//transientDataMap["result"] = []byte("Transient data in save participant")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID2, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID2, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(args[5]), []byte(args[6])}, TransientMap: transientDataMap})
	if err != nil {
		fmt.Errorf("failed to make external payment: %v", err)
		return "", err
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
		return "", errors.New("did NOT receive CC event for eventId")
	}

	return string(response.TransactionID), nil

}

// WalletPayment driver
func (setup *FabricSetup) WalletPayment(walletPayment Models.WalletPayment) (string, error) {
	// Prepare arguments
	var args []string
	args = append(args, "makeWalletPayment")
	args = append(args, walletPayment.FromWallet)
	args = append(args, walletPayment.ToWallet)
	args = append(args, walletPayment.Date)
	args = append(args, strconv.Itoa(walletPayment.Quantity))
	args = append(args, externalPayment.Identification)

	eventID := "makeWalletPayment"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	//transientDataMap["result"] = []byte("Transient data in save participant")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID2, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID2, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(args[5]), []byte(args[6])}, TransientMap: transientDataMap})
	if err != nil {
		fmt.Errorf("failed to make external payment: %v", err)
		return "", err
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
		return "", errors.New("did NOT receive CC event for eventId")
	}

	return string(response.TransactionID), nil

}
