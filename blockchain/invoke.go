package blockchain

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/pkg/errors"

	Models "github.com/sumaikun/apeslogistic-rest-api/models"
)

// InvokeHello driver
func (setup *FabricSetup) InvokeHello(value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "hello")
	args = append(args, value)

	eventID := "eventInvoke"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}, TransientMap: transientDataMap})
	if err != nil {
		fmt.Errorf("failed to move funds: %v", err)
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

// SaveParticipant driver
func (setup *FabricSetup) SaveParticipant(participant Models.Participant) (string, error) {
	// Prepare arguments
	var args []string
	args = append(args, "saveParticipant")
	args = append(args, participant.Name)
	args = append(args, participant.Nationality)
	args = append(args, participant.Address)
	args = append(args, participant.Phone)
	args = append(args, participant.Identification)
	args = append(args, participant.PhotoURL)
	args = append(args, participant.Description)

	eventID := "saveParticipant"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in create participant")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(args[5]), []byte(args[6]), []byte(args[7])}, TransientMap: transientDataMap})
	if err != nil {
		fmt.Errorf("failed to move funds: %v", err)
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

// SaveAsset driver
func (setup *FabricSetup) SaveAsset(asset Models.Asset) (string, error) {
	// Prepare arguments
	var args []string
	args = append(args, "saveAsset")
	args = append(args, asset.Participant)
	args = append(args, asset.State)
	args = append(args, asset.Location)
	args = append(args, asset.Meta)
	args = append(args, asset.Identification)
	args = append(args, asset.Title)
	args = append(args, asset.Date)
	args = append(args, asset.AssetType)

	eventID := "saveAsset"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in create asset")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(args[5]), []byte(args[6]), []byte(args[7]), []byte(args[8])}, TransientMap: transientDataMap})
	if err != nil {
		fmt.Errorf("failed to move funds: %v", err)
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

//--------------------------------------------------- wallet driver -------------------------------------------------------------------//

// SaveOwner driver
func (setup *FabricSetup) SaveOwner(owner Models.Owner) (string, error) {
	// Prepare arguments
	var args []string
	args = append(args, "saveOwner")
	args = append(args, owner.Name)
	args = append(args, owner.Nationality)
	args = append(args, owner.Address)
	args = append(args, owner.Phone)
	args = append(args, owner.Identification)
	args = append(args, owner.PhotoURL)
	args = append(args, owner.Notes)

	eventID := "saveOwner"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	//transientDataMap["result"] = []byte("Transient data in save participant")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID2, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(args[5]), []byte(args[6]), []byte(args[7])}, TransientMap: transientDataMap})
	if err != nil {
		fmt.Errorf("failed to save owner: %v", err)
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

// SaveExternalAgent driver
func (setup *FabricSetup) SaveExternalAgent(externalAgent Models.ExternalAgent) (string, error) {
	// Prepare arguments
	var args []string
	args = append(args, "saveExternalAgent")
	args = append(args, externalAgent.Name)
	args = append(args, externalAgent.Description)
	args = append(args, externalAgent.Identification)

	eventID := "saveExternalAgent"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	//transientDataMap["result"] = []byte("Transient data in save participant")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID2, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}, TransientMap: transientDataMap})
	if err != nil {
		fmt.Errorf("failed to save external agent: %v", err)
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

// SaveEvent driver
func (setup *FabricSetup) SaveEvent(event Models.Event) (string, error) {
	// Prepare arguments
	var args []string
	args = append(args, "saveEvent")
	args = append(args, event.FromExternal)
	args = append(args, event.FromWallet)
	args = append(args, event.ToWallet)
	args = append(args, event.ToExternal)

	eventID := "saveEvent"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	//transientDataMap["result"] = []byte("Transient data in save participant")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID2, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4])}, TransientMap: transientDataMap})
	if err != nil {
		fmt.Errorf("failed to save external agent: %v", err)
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

// SaveRule driver
func (setup *FabricSetup) SaveRule(rule Models.Rule) (string, error) {
	// Prepare arguments
	var args []string
	args = append(args, "saveRule")
	args = append(args, rule.Event)
	args = append(args, rule.Fee)
	args = append(args, rule.ToWallet)
	args = append(args, rule.ToExternal)
	args = append(args, rule.Date)
	args = append(args, rule.Quantity)
	args = append(args, rule.State)

	eventID := "saveRule"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	//transientDataMap["result"] = []byte("Transient data in save participant")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID2, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(args[5]), []byte(args[6]), []byte(args[7])}, TransientMap: transientDataMap})
	if err != nil {
		fmt.Errorf("failed to save external agent: %v", err)
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
