package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/manifoldco/go-signature"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type Database struct {
	Resources   map[string]string `json:"resources"`
	Credentials map[string]string `json:"credentials"`
}

type ResponseStruct struct {
	Message string `json:"message"`
}

type Resources struct {
	Id           string `json:"id"`
	Product      string `json:"product"`
	Plan         string `json:"plan"`
	Region       string `json:"region"`
	RandomNumber string `json:"randomNumber"`
}

// @TODO: may want to rename these
type CredentialsResponse struct {
	Message     string `json:"message"`
	Credentials `json:"credentials"`
}

type Credentials struct {
	Password string `json:"password"`
}

type CredentialsRequest struct {
	Id         string `json:"id"`
	ResourceId string `json:"resource_id"`
}

// @TODO: refactor following "get..." functions
func GetBodyBufferAndResources(r *http.Request) (*bytes.Buffer, Resources) {
	body, _ := ioutil.ReadAll(r.Body)
	bodyBuffer := bytes.NewBuffer(body)
	bodyCopy := bytes.NewReader(body) // clone body to avoid mutability issues
	rqs := GetResources(bodyCopy)

	return bodyBuffer, rqs
}

func GetResources(bodyCopy io.Reader) Resources {
	// decode request body
	decoder := json.NewDecoder(bodyCopy)
	var rqs Resources
	decoder.Decode(&rqs)

	return rqs
}

func GetBodyBufferAndCredentials(r *http.Request) (*bytes.Buffer, CredentialsRequest) {
	body, _ := ioutil.ReadAll(r.Body)
	bodyBuffer := bytes.NewBuffer(body)
	bodyCopy := bytes.NewReader(body) // clone body to avoid mutability issues
	rqs := GetCredentials(bodyCopy)

	return bodyBuffer, rqs
}

func GetCredentials(bodyCopy io.Reader) CredentialsRequest {
	// decode request body
	decoder := json.NewDecoder(bodyCopy)
	var rqs CredentialsRequest
	decoder.Decode(&rqs)

	return rqs
}

func SignatureIsNotValid(r *http.Request, w http.ResponseWriter, buf io.Reader) bool {
	verifier, _ := signature.NewVerifier(MASTER_KEY)

	if err := verifier.Verify(r, buf); err != nil {
		resp := ResponseStruct{""}
		js, err := json.Marshal(resp)

		IssueResponseIfErrorOccurs(err, w)

		w.WriteHeader(http.StatusUnauthorized)
		w.Write(js)
		return true
	}

	return false
}

func Seed() string {
	rand.Seed(time.Now().UTC().UnixNano())
	result := rand.Intn(100)
	return fmt.Sprintf("%d", result)
}

func HandleResponse(responseMessage string, statusCode int, w http.ResponseWriter) {
	resp := ResponseStruct{responseMessage}
	js, err := json.Marshal(resp)

	IssueResponseIfErrorOccurs(err, w)

	w.WriteHeader(statusCode)
	w.Write(js)
}

func ConvertRequestToJson(rqs Resources) []byte {
	rqsData := Resources{rqs.Id, rqs.Product, rqs.Plan, rqs.Region, rqs.RandomNumber}
	jsonData, _ := json.Marshal(rqsData)

	return jsonData
}

func IssueResponseIfErrorOccurs(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
