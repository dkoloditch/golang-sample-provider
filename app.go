package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	// Manifold's public key (or your local test version), used
	// for by the manifoldco_signature package to verifiy that requests came
	// from Manifold.
	MASTER_KEY    = os.Getenv("MASTER_KEY")

	// OAuth 2.0 client id and secret pair. Used to exchange a code for a user's
	// token during SSO.
	CLIENT_ID     = os.Getenv("CLIENT_ID")
	CLIENT_SECRET = os.Getenv("CLIENT_SECRET")

	// The URL of manifold's connector url, for completing SSO or making requests
	CONNECTOR_URL = os.Getenv("CONNECTOR_URL")

	// oauth2 config setup for SSO
	oac           = &oauth2.Config{
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  CONNECTOR_URL,
			TokenURL: CONNECTOR_URL + "/oauth/tokens",
		},
	}

	// Products, plans, and regions we know about
	products = [1]string{"bonnets"}
	plans = [2]string{"small", "large"}
	regions = [1]string{"aws::us-east-1"}

	// simple in-memory db using structs and nested maps
	db = Database{
		Resources:   make(map[string]string),
		Credentials: make(map[string]string),
	}
)

func main() {
	// Manifold API endpoints and functions for handling requests
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/dashboard", dashboardHandler).Methods("GET")

	router.HandleFunc("/v1/resources/{id}", createResourceHandler).Methods("PUT")
	router.HandleFunc("/v1/resources/{id}", updateResourceHandler).Methods("PATCH")
	router.HandleFunc("/v1/resources/{id}", deleteResourceHandler).Methods("DELETE")

	router.HandleFunc("/v1/credentials/{id}", createCredentialHandler).Methods("PUT")
	router.HandleFunc("/v1/credentials/{id}", deleteCredentialHandler).Methods("DELETE")

	router.HandleFunc("/v1/sso", ssoHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":4567", router))
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// The cool dashboard. A user has to be authenticated with Manifold to
	// use this.

	return
}

func createResourceHandler(w http.ResponseWriter, r *http.Request) {
	SetContentTypeHeaderAsJSON(w)

	bodyBuffer, rqs := GetBodyBufferAndResources(r)
	id := rqs.Id

	if SignatureIsNotValid(r, w, bodyBuffer) {
		return
	}

	if ProductIsNotValid(rqs.Product, w) {
		return
	}

	if PlanIsNotValid(rqs.Plan, w) {
		return
	}

	if RegionIsNotValid(rqs.Region, w) {
		return
	}

	if ResourceAlreadyExists(rqs, w, id) {
		return
	}

	if ResourceCreated(rqs, w) {
		return
	}
}

func updateResourceHandler(w http.ResponseWriter, r *http.Request) {
	SetContentTypeHeaderAsJSON(w)

	// since the id is only passed via URL with PATCH requests, we set this here
	// and provide it to the relevant methods below.
	_, id := path.Split(r.URL.Path)
	bodyBuffer, rqs := GetBodyBufferAndResources(r)

	if SignatureIsNotValid(r, w, bodyBuffer) {
		return
	}

	if PlanIsNotValid(rqs.Plan, w) {
		return
	}

	if ResourceDoesNotExist(rqs, w, id) {
		return
	}

	if ResourceUpdated(rqs, w, id) {
		return
	}
}

func deleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	SetContentTypeHeaderAsJSON(w)

	_, id := path.Split(r.URL.Path)
	bodyBuffer, rqs := GetBodyBufferAndResources(r)

	if SignatureIsNotValid(r, w, bodyBuffer) {
		return
	}

	if ResourceDoesNotExist(rqs, w, id) {
		return
	}

	if ResourceDeleted(w, id) {
		return
	}

	return
}

func createCredentialHandler(w http.ResponseWriter, r *http.Request) {
	SetContentTypeHeaderAsJSON(w)

	bodyBuffer, rqs := GetBodyBufferAndCredentials(r)

	if SignatureIsNotValid(r, w, bodyBuffer) {
		return
	}

	if InvalidResourceId(w, rqs) {
		return
	}

	if ProvisionCredentials(w, rqs) {
		return
	}

	return
}

func deleteCredentialHandler(w http.ResponseWriter, r *http.Request) {
	SetContentTypeHeaderAsJSON(w)

	_, id := path.Split(r.URL.Path)

	if CredentialsDoNotExist(w, id) {
		return
	}

	if CredentialsDeleted(w, id) {
		return
	}
}

func ssoHandler(w http.ResponseWriter, r *http.Request) {
	// SSO requets come from the user's browser, so we don't want to run the
	// validator against them.

	SetContentTypeHeaderAsJSON(w)

	code := r.URL.Query().Get("code")
	// resource_id := r.URL.Query().Get("resource_id")
	// url := oac.AuthCodeURL(CONNECTOR_URL)
	_, err := oac.Exchange(oauth2.NoContext, code)
	errReformatted := fmt.Errorf("%v", err) // avoids a blowup

	if err != nil {
		resp := &CredentialsResponse{
			Message: fmt.Sprintf("%v", errReformatted.Error()),
		}
		js, err := json.Marshal(resp)

		IssueResponseIfErrorOccurs(err, w)

		w.WriteHeader(http.StatusUnauthorized)
		w.Write(js)
	}

	// @TODO: store token and resource to session in db
	// @TODO: redirect

	return
}
