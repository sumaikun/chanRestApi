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

//------------------------------- wallet chaincode -------------------------

// Owner representation in chaincode
type Owner struct {
	ObjectType     string `json:"docType"`
	Name           string `json:"name"`
	Nationality    string `json:"nationality"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Identification string `json:"identification"`
	PhotoURL       string `json:"photoUrl"`
	Notes          string `json:"notes"`
	Balance        int    `json:"balance"`
}

// ExternalAgent representation in chaincode
type ExternalAgent struct {
	ObjectType     string `json:"docType"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Identification string `json:"identification"`
}

// WalletPayment representation in chaincode
type WalletPayment struct {
	ObjectType     string `json:"docType"`
	FromWallet     string `json:"fromWallet"`
	ToWallet       string `json:"toWallet"`
	State          string `json:"state"`
	Date           string `json:"date"`
	Quantity       int    `json:"quantity"`
	Identification string `json:"identification"`
}

// ExternalPayment representation in chaincode
type ExternalPayment struct {
	ObjectType     string `json:"docType"`
	FromExternal   string `json:"fromExternal"`
	ToWallet       string `json:"toWallet"`
	State          string `json:"state"`
	Date           string `json:"date"`
	Quantity       int    `json:"quantity"`
	PaymentType    string `json:"paymentType"`
	Identification string `json:"identification"`
}

// Event representation in chaincode
type Event struct {
	ObjectType     string `json:"docType"`
	FromExternal   string `json:"fromExternal"`
	FromWallet     string `json:"fromWallet"`
	ToWallet       string `json:"toWallet"`
	ToExternal     string `json:"toExternal"`
	Identification string `json:"identification"`
}
