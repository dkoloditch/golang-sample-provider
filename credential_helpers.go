package main

import (
	// "fmt"
	"encoding/json"
	"net/http"
)

func ProvisionCredentials(w http.ResponseWriter, rqs CredentialsRequest) bool {
	data, err := json.Marshal(rqs)

	IssueResponseIfErrorOccurs(err, w)

	db.Credentials[rqs.Id] = string(data)

	message := "your password is ready"

	setHeadersAndResponse(message, http.StatusCreated, RESPONSE_TYPE_CREDENTIAL, w)

	return true
}

func CredentialsDoNotExist(w http.ResponseWriter, id string) bool {
	_, dataRetrieved := db.Credentials[id]

	if !dataRetrieved {
		message := "no such credential"

		setHeadersAndResponse(message, http.StatusNotFound, RESPONSE_TYPE_CREDENTIAL, w)

		return true
	}

	return false
}

func CredentialsDeleted(w http.ResponseWriter, id string) bool {
	delete(db.Credentials, id)

	message := ""

	setHeadersAndResponse(message, http.StatusNoContent, RESPONSE_TYPE_CREDENTIAL, w)

	return true
}
