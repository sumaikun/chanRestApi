package main

import (
	//"fmt"
	//"net/http"
	//"net/url"
	//"time"

	"net/http"

	Models "github.com/sumaikun/apeslogistic-rest-api/models"
	"github.com/thedevsaddam/govalidator"
	//"gopkg.in/mgo.v2/bson"
)

func userValidator(r *http.Request) (map[string]interface{}, Models.User) {

	var user Models.User

	rules := govalidator.MapData{
		"name":    []string{"required"},
		"email":   []string{"required", "email"},
		"phone":   []string{"min:7", "max:10"},
		"address": []string{"required"},
		//"picture": []string{"url"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &user,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, user
}

func participantValidator(r *http.Request) (map[string]interface{}, Models.Participant) {

	var participant Models.Participant

	rules := govalidator.MapData{
		"name":           []string{"required"},
		"nationality":    []string{"required"},
		"address":        []string{"required"},
		"phone":          []string{"required"},
		"identification": []string{"required"},
		"description":    []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &participant,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, participant
}

func assetValidator(r *http.Request) (map[string]interface{}, Models.Asset) {

	var asset Models.Asset

	rules := govalidator.MapData{
		"participant":    []string{"required"},
		"state":          []string{"required"},
		"location":       []string{"required"},
		"title":          []string{"required"},
		"identification": []string{"required"},
		"assetType":      []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &asset,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, asset
}
