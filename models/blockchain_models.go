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
