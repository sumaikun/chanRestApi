package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	Helpers "github.com/sumaikun/apeslogistic-rest-api/helpers"
	Models "github.com/sumaikun/apeslogistic-rest-api/models"
	"gopkg.in/mgo.v2/bson"
)

//-----------------------------  Auth functions --------------------------------------------------

func authentication(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	response := &Models.TokenResponse{Token: "", User: nil}

	var creds Models.Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := Models.Users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || !Helpers.CheckPasswordHash(creds.Password, expectedPassword) {

		user, err := dao.FindOneByKEY("users", "email", creds.Username)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		match := Helpers.CheckPasswordHash(creds.Password, user.(bson.M)["password"].(string))

		if !match {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		response.User = user.(bson.M)

	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(8 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	w.Header().Set("Content-type", "application/json")

	//Generate json response for get the token
	response.Token = tokenString

	json.NewEncoder(w).Encode(response)
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"status":"ok"}`)
}

func createInititalUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	users, err := dao.FindAll("users")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if users == nil {

		fmt.Println("is nil")

		var user Models.User

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user.ID = bson.NewObjectId()
		user.Date = time.Now().String()
		user.UpdateDate = time.Now().String()

		if len(user.Password) != 0 {
			user.Password, _ = Helpers.HashPassword(user.Password)
		}

		if err := dao.Insert("users", user, []string{"email"}); err != nil {
			Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		Helpers.RespondWithJSON(w, http.StatusCreated, user)

	} else {
		Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "can not create initial users again"})
	}

}

//-----------------------------  Users functions --------------------------------------------------

func allUsersEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	users, err := dao.FindAll("users")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, users)
}

func createUsersEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	usera := context.Get(r, "user")

	userParsed := usera.(bson.M)

	w.Header().Set("Content-type", "application/json")

	err, user := userValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	user.ID = bson.NewObjectId()
	user.Date = time.Now().String()
	user.UpdateDate = time.Now().String()
	user.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	user.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if len(user.Password) != 0 {
		user.Password, _ = Helpers.HashPassword(user.Password)
	}

	if err := dao.Insert("users", user, []string{"email"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, user)

}

func findUserEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	user, err := dao.FindByID("users", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, user)

}

func removeUserEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("users", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateUserEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	params := mux.Vars(r)

	usera := context.Get(r, "user")

	userParsed := usera.(bson.M)

	w.Header().Set("Content-type", "application/json")

	err, user := userValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevUser, err2 := dao.FindByID("users", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	parsedData := prevUser.(bson.M)

	user.ID = parsedData["_id"].(bson.ObjectId)

	user.Date = parsedData["date"].(string)

	user.UpdateDate = time.Now().String()

	if parsedData["createdBy"] == nil {
		user.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	} else {
		user.CreatedBy = parsedData["createdBy"].(string)
	}

	user.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if len(user.Password) == 0 {
		user.Password = parsedData["password"].(string)
	} else {
		user.Password, _ = Helpers.HashPassword(user.Password)
	}

	if err := dao.Update("users", user.ID, user); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//-------------------------------------- file Upload -----------------------------------------

func fileUpload(w http.ResponseWriter, r *http.Request) {

	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	var extension = filepath.Ext(handler.Filename)

	fmt.Printf("Extension: %+v\n", extension)

	tempFile, err := ioutil.TempFile("files", "upload-*"+extension)

	if err != nil {
		fmt.Println(err)
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
	}

	var tempPath = tempFile.Name()

	fmt.Println("temp file before trim" + tempPath)

	var tempName = strings.Replace(tempPath, "files/", "", -1)

	fmt.Println("tempName " + tempName)

	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"filename": tempName})

}

func deleteImage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var fileName = params["file"]

	var err = os.Remove("./files/" + fileName)
	if err != nil {
		//log.Fatal(err) // perhaps handle this nicer
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "fileDelete"})
	return

}

func serveImage(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var fileName = params["image"]

	if !strings.Contains(fileName, "png") && !strings.Contains(fileName, "jpg") && !strings.Contains(fileName, "jpeg") && !strings.Contains(fileName, "gif") {
		Helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "invalid file extension"})
		return
	}

	img, err := os.Open("./files/" + params["image"])
	if err != nil {
		//log.Fatal(err) // perhaps handle this nicer
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	defer img.Close()
	w.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
	io.Copy(w, img)

}

func downloadFile(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var fileName = params["file"]

	http.ServeFile(w, r, "./files/"+fileName)
}

func getFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

//-----------------------------  Hello chaincode --------------------------------------------------

func (app *Application) queryHelloChainCode(w http.ResponseWriter, r *http.Request) {

	// Query the chaincode
	response, err := app.Fabric.QueryHello()
	if err != nil {
		fmt.Printf("Unable to query hello on the chaincode: %v\n", err)
		Helpers.RespondWithJSON(w, http.StatusBadGateway, err)
		return
	}

	fmt.Printf("Response from the query hello: %s\n", response)
	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": response})

}

func (app *Application) invokeHelloChaincode(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	if len(params["word"]) == 0 {
		Helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "param word needed"})
	}

	defer r.Body.Close()

	// Invoke the chaincode
	txID, err := app.Fabric.InvokeHello(params["word"])
	if err != nil {
		fmt.Printf("Unable to invoke hello on the chaincode: %v\n", err)
		Helpers.RespondWithJSON(w, http.StatusBadGateway, err)
		return
	}
	fmt.Printf("Successfully invoke hello, transaction ID: %s\n", txID)
	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": txID})

}

func (app *Application) historyHelloChainCode(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Invoke the chaincode
	response, err := app.Fabric.HistoryHello()
	if err != nil {
		fmt.Printf("Unable to history hello on the chaincode: %v\n", err)
		Helpers.RespondWithJSON(w, http.StatusBadGateway, err)
		return
	}

	fmt.Printf("Response from the history hello: %s\n", response)

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": response})

}

func (app *Application) getDataFromChaincode(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	if len(params["key"]) == 0 {
		Helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "param key needed"})
	}

	defer r.Body.Close()

	response, err := app.Fabric.QueryGetData(params["key"])
	if err != nil {
		fmt.Printf("Unable to query  the chaincode: %v\n", err)
		Helpers.RespondWithJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	fmt.Printf("Response from chaincode: %s\n", response)

	/*out, err := json.Marshal(response)
	if err != nil {
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}*/

	Helpers.RespondWithJSON(w, http.StatusOK, response)

}

func (app *Application) saveParticipant(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	err, participant := participantValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	// Invoke the chaincode
	txID, err2 := app.Fabric.SaveParticipant(participant)
	if err2 != nil {
		fmt.Printf("Unable to save participant on the chaincode: %v\n", err2)
		Helpers.RespondWithJSON(w, http.StatusBadGateway, map[string]string{"error": err2.Error()})
		return
	}
	fmt.Printf("Successfully save participant transaction ID: %s\n", txID)
	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": txID})

}

func (app *Application) getParticipants(w http.ResponseWriter, r *http.Request) {

	response, err := app.Fabric.QueryObjectType("participant")
	if err != nil {
		fmt.Printf("Unable to query  the chaincode: %v\n", err)
		Helpers.RespondWithJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	fmt.Printf("Response from chaincode: %s\n", response)
	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": response})

}

func (app *Application) saveAsset(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	err, asset := assetValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	asset.Date = time.Now().String()

	// Invoke the chaincode
	txID, err2 := app.Fabric.SaveAsset(asset)
	if err2 != nil {
		fmt.Printf("Unable to save asset on the chaincode: %v\n", err2)
		Helpers.RespondWithJSON(w, http.StatusBadGateway, map[string]string{"error": err2.Error()})
		return
	}
	fmt.Printf("Successfully save asset transaction ID: %s\n", txID)
	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": txID})

}

func (app *Application) getAssets(w http.ResponseWriter, r *http.Request) {

	response, err := app.Fabric.QueryObjectType("asset")
	if err != nil {
		fmt.Printf("Unable to query  the chaincode: %v\n", err)
		Helpers.RespondWithJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	fmt.Printf("Response from chaincode: %s\n", response)
	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": response})

}
