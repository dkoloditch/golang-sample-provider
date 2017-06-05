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

	resp := &CredentialsResponse{
		Message: "your password is ready",
		Credentials: Credentials{
			Password: "test1234",
		},
	}
	js, err := json.Marshal(resp)

	IssueResponseIfErrorOccurs(err, w)

	w.WriteHeader(http.StatusCreated)
	w.Write(js)

	return true
}

func CredentialsDoNotExist(w http.ResponseWriter, id string) bool {
	_, dataRetrieved := db.Credentials[id]

	if !dataRetrieved {
		resp := &CredentialsResponse{
			Message: "no such credential",
		}
		js, err := json.Marshal(resp)

		IssueResponseIfErrorOccurs(err, w)

		w.WriteHeader(http.StatusNotFound)
		w.Write(js)

		return true
	}

	return false
}

func CredentialsDeleted(w http.ResponseWriter, id string) bool {
	delete(db.Credentials, id)

	resp := &CredentialsResponse{}
	js, err := json.Marshal(resp)

	IssueResponseIfErrorOccurs(err, w)

	w.WriteHeader(http.StatusNoContent)
	w.Write(js)

	return true
}
