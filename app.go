package main

import (
  "fmt"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  // "golang.org/x/oauth2"
  "os"
  "math/rand"
  "time"
  "encoding/json"
  "github.com/manifoldco/go-signature"
  "io"
  "io/ioutil"
  "bytes"
)

// simple in-memory db
var db = make(map[string]string)

var MASTER_KEY = os.Getenv("MASTER_KEY")
// CLIENT_ID := os.Getenv("CLIENT_ID")
// CLIENT_SECRET := os.Getenv("CLIENT_SECRET")
// CONNECTOR_URL := os.Getenv("CONNECTOR_URL")

type ResponseStruct struct {
  Message string `json:"message"`
}

// @TODO: may want to rename this
type RequestStruct struct {
  Id string
  Product string
  Plan string
  Region string
}

func main() {
  router := mux.NewRouter().StrictSlash(true)

  router.HandleFunc("/dashboard", dashboardHandler).Methods("GET")

  router.HandleFunc("/v1/resources/{id}", createResourcesHandler).Methods("PUT")
  router.HandleFunc("/v1/resources/{id}", updateResourcesHandler).Methods("PATCH")
  router.HandleFunc("/v1/resources/{id}", deleteResourcesHandler).Methods("DELETE")

  router.HandleFunc("/v1/credentials/{id}", createCredentialsHandler).Methods("PUT")
  router.HandleFunc("/v1/credentials/{id}", deleteCredentialsHandler).Methods("DELETE")

  router.HandleFunc("/v1/sso", ssoHandler).Methods("GET")

  log.Fatal(http.ListenAndServe(":4567", router))
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
  return
}

func createResourcesHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  bodyBuffer, rqs := getBodyBufferAndRequestStruct(r)

  if signatureIsNotValidAndResponseCreated(r, w, bodyBuffer) { return }

  if productIsNotValidAndResponseCreated(rqs.Product, w) { return }

  if planIsNotValidAndResponseCreated(rqs.Plan, w) { return }

  if regionIsNotValidAndResponseCreated(rqs.Region, w) { return }

  if resourceAlreadyExistsAndResponseCreated(rqs, w) { return }

  if validRequestAndResponseCreated(rqs, w) { return }
}

func updateResourcesHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  bodyBuffer, rqs := getBodyBufferAndRequestStruct(r)

  if signatureIsNotValidAndResponseCreated(r, w, bodyBuffer) { return }

  if resourceAlreadyExistsAndResponseCreated(rqs, w) { return }

  if resourceDoesNotExistAndResponseCreated(rqs, w) { return }
}

func deleteResourcesHandler(w http.ResponseWriter, r *http.Request) {
  return
}

func createCredentialsHandler(w http.ResponseWriter, r *http.Request) {
  return
}

func deleteCredentialsHandler(w http.ResponseWriter, r *http.Request) {
  return
}

func ssoHandler(w http.ResponseWriter, r *http.Request) {
  return
}

func getBodyBufferAndRequestStruct(r *http.Request) (*bytes.Buffer, RequestStruct) {
  body, _ := ioutil.ReadAll(r.Body)
  bodyBuffer := bytes.NewBuffer(body)
  bodyCopy := bytes.NewReader(body) // clone body to avoid mutability issues
  rqs := getRequestStruct(bodyCopy)

  return bodyBuffer, rqs
}

func getRequestStruct(bodyCopy io.Reader) RequestStruct {
  // decode request body
  decoder := json.NewDecoder(bodyCopy)
  var rqs RequestStruct
  decoder.Decode(&rqs)

  return rqs
}

func signatureIsNotValidAndResponseCreated(r *http.Request, w http.ResponseWriter, buf io.Reader) bool {
  verifier, _ := signature.NewVerifier(MASTER_KEY)

  if err := verifier.Verify(r, buf); err != nil {
    resp := ResponseStruct{""}
    js, err := json.Marshal(resp)

    issueResponseIfErrorOccurs(err, w)

    w.WriteHeader(http.StatusUnauthorized)
    w.Write(js)
    return true
  }

  return false
}

func seed() string {
	rand.Seed(time.Now().UTC().UnixNano())
	result := rand.Intn(100)
  return fmt.Sprintf("%d", result)
}

func productIsNotValidAndResponseCreated(product string, w http.ResponseWriter) bool {
  products := [1]string{"bonnets"}

  for i := range products {
    if product == products[i] {
      return false
    }
  }

  handleResponse("bad product", http.StatusBadRequest, w)

  return true
}

func planIsNotValidAndResponseCreated(plan string, w http.ResponseWriter) bool {
  plans := [2]string{"small", "large"}

  for i := range plans {
    if plan == plans[i] {
      return false
    }
  }

  handleResponse("bad plan", http.StatusBadRequest, w)

  return true
}

func regionIsNotValidAndResponseCreated(region string, w http.ResponseWriter) bool {
  regions  := [1]string{"aws::us-east-1"}

  for i := range regions {
    if region == regions[i] {
      return false
    }
  }

  handleResponse("bad region", http.StatusBadRequest, w)

  return true
}

func resourceAlreadyExistsAndResponseCreated(rqs RequestStruct, w http.ResponseWriter) bool {
  existingData, dataRetrieved := db[rqs.Id]
  newData := string(convertRequestToJson(rqs))
  newDataMatchesOldData := existingData == newData

  // @TODO this can probably be refactored to not account for Method given
  // appropriate routes
  if dataRetrieved {
    // same content acts as created
    if newDataMatchesOldData {
      resp := ResponseStruct{""}
      js, err := json.Marshal(resp)

      issueResponseIfErrorOccurs(err, w)

      w.WriteHeader(http.StatusCreated)
      w.Write(js)
    } else {
      // different content results in conflict
      resp := ResponseStruct{"resource already exists"}
      js, err := json.Marshal(resp)

      issueResponseIfErrorOccurs(err, w)

      w.WriteHeader(http.StatusConflict)
      w.Write(js)
    }

    return true
  }

  return false
}

func resourceDoesNotExistAndResponseCreated(rqs RequestStruct, w http.ResponseWriter) bool {
  _, dataRetrieved := db[rqs.Id]

  if !dataRetrieved {
    // non existing resource
    resp := ResponseStruct{"no such resource"}
    js, err := json.Marshal(resp)

    issueResponseIfErrorOccurs(err, w)

    w.WriteHeader(http.StatusNotFound)
    w.Write(js)

    return true
  }

  return false
}

func validRequestAndResponseCreated(rqs RequestStruct, w http.ResponseWriter) bool {
  // @TODO: can this be abstracted further?

  // get the random number and create json response
  result := seed()
  resp := ResponseStruct{fmt.Sprintf("%d", result)}
  rqsData := RequestStruct{rqs.Id, rqs.Product, rqs.Plan, rqs.Region}
  js, err := json.Marshal(resp)
  data, err := json.Marshal(rqsData)

  issueResponseIfErrorOccurs(err, w)

  // add to db
  db[rqs.Id] = string(data)

  w.WriteHeader(http.StatusCreated)
  w.Write(js)

  return true
}

func handleResponse(responseMessage string, statusCode int, w http.ResponseWriter) {
  resp := ResponseStruct{responseMessage}
  js, err := json.Marshal(resp)

  issueResponseIfErrorOccurs(err, w)

  w.WriteHeader(statusCode)
  w.Write(js)
}

func convertRequestToJson(rqs RequestStruct) []byte {
  rqsData := RequestStruct{rqs.Id, rqs.Product, rqs.Plan, rqs.Region}
  jsonData, _ := json.Marshal(rqsData)

  return jsonData
}

func issueResponseIfErrorOccurs(err error, w http.ResponseWriter) {
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}
