package models

import "gopkg.in/mgo.v2/bson"

//User representation on mongo
type User struct {
	ID             bson.ObjectId `bson:"_id" json:"id"`
	Name           string        `bson:"name" json:"name"`
	LastName       string        `bson:"lastName" json:"lastName"`
	Password       string        `bson:"password" json:"password"`
	Email          string        `bson:"email" json:"email"`
	Address        string        `bson:"address" json:"address"`
	Role           string        `bson:"role" json:"role"`
	Phone          string        `bson:"phone" json:"phone"`
	Picture        string        `bson:"picture" json:"picture"`
	State          string        `bson:"state" json:"state"`
	DocumentType   string        `bson:"documentType" json:"documentType"`
	DocumentNumber string        `bson:"documentNumber" json:"documentNumber"`
	Participants   []string      `bson:"participants" json:"participants"`
	CreatedBy      string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy      string        `bson:"updatedBy" json:"updatedBy"`
	Date           string        `bson:"date" json:"date"`
	UpdateDate     string        `bson:"update_date" json:"update_date"`
}
