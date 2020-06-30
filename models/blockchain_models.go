package models

//Participant representation on fabric
type Participant struct {
	ID                   string `json:"id"`
	DisplayName          string `json:"displayName"`
	Email                string `json:"email"`
	ServiceType          string `json:"serviceType"`
	ExporterConfirmation bool   `json:"exporterConfirmation"`
	IntegrationLevel     string `json:"integrationLevel"`
}

//Issue or token type representation on fabric
type Issue struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Participant string `json:"participant"`
	Description string `json:"description"`
}

//Assets representation on fabric
type Assets struct {
	ID    string `json:"id"`
	Issue string `json:"issue"`
	Meta  string `json:"meta"`
}
