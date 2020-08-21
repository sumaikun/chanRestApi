package models

//Participant representation on fabric
type Participant struct {
	Name           string `json:"name"`
	Nationality    string `json:"nationality"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	Identification string `json:"identification"`
	PhotoURL       string `json:"photoUrl"`
	Description    string `json:"description"`
}

// Asset representation in chaincode
type Asset struct {
	ObjectType     string `json:"docType"`
	Participant    string `json:"participant"`
	State          string `json:"state"`
	Location       string `json:"location"`
	Meta           string `json:"meta"`
	Identification string `json:"identification"`
	Title          string `json:"title"`
	Date           string `json:"date"`
	AssetType      string `json:"assetType"`
}
