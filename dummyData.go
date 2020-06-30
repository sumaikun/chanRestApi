package main

import (
	"net/http"

	Models "github.com/sumaikun/apeslogistic-rest-api/models"
)

type allParticipants []Models.Participant

var participants = allParticipants{
	{
		ID:                   "DSDSADA12122",
		DisplayName:          "Becket",
		Email:                "info@becket.com",
		ServiceType:          "fullfillment",
		ExporterConfirmation: true,
		IntegrationLevel:     "Full Rfid",
	},
	{
		ID:                   "DSDSADA12142",
		DisplayName:          "Unibán",
		Email:                "info@uniban.com",
		ServiceType:          "control and trazability",
		ExporterConfirmation: true,
		IntegrationLevel:     "Full Rfid",
	},
}

type allIssues []Models.Issue

var issues = allIssues{
	{
		ID:          "32123ddsadsad2",
		Name:        "Cinturones",
		Participant: "Becket",
		Description: "Cinturones fabricados por becket colombia",
	},
	{
		ID:          "32123ddsafsad2",
		Name:        "Zapatos",
		Participant: "Unibán",
		Description: "Zapatos fabricados por becket colombia",
	},
	{
		ID:          "32153ddsafsad2",
		Name:        "Bananos",
		Participant: "Unibán",
		Description: "Activo de cargamento de fruta Unibán",
	},
}

type allAssets []Models.Assets

var assets = allAssets{}

func getParticipants(w http.ResponseWriter, r *http.Request) {

}

func getAssets(w http.ResponseWriter, r *http.Request) {

}

func getIssues(w http.ResponseWriter, r *http.Request) {

}
