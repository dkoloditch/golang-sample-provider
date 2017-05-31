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

type RequestStruct struct {
  Id string
  Product string
  Plan string
  Region string
}

func main() {
  router := mux.NewRouter().StrictSlash(true)

  router.HandleFunc("/dashboard", dashboardHandler).Methods("GET")

  router.HandleFunc("/v1/resources/{id}", resourcesHandler).Methods("PUT")
  router.HandleFunc("/v1/resources/{id}", resourcesHandler).Methods("PATCH")
  router.HandleFunc("/v1/resources/{id}", resourcesHandler).Methods("DELETE")

  router.HandleFunc("/v1/credentials/{id}", credentialsHandler).Methods("PUT")
  router.HandleFunc("/v1/credentials/{id}", credentialsHandler).Methods("DELETE")

  router.HandleFunc("/v1/sso", ssoHandler).Methods("GET")

  log.Fatal(http.ListenAndServe(":4567", router))
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
  return
}

func resourcesHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  body, _ := ioutil.ReadAll(r.Body)
  buf := bytes.NewBuffer(body)
  bodyCopy := bytes.NewReader(body) // clone body to avoid mutability issues
  verifier, _ := signature.NewVerifier(MASTER_KEY)
  if err := verifier.Verify(r, buf); err != nil {
    resp := ResponseStruct{""}
    js, err := json.Marshal(resp)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    w.WriteHeader(http.StatusUnauthorized)
    w.Write(js)
    return
  }

  // decode request body
  decoder := json.NewDecoder(bodyCopy)
  var rqs RequestStruct
  decoder.Decode(&rqs)

  if !checkProduct(rqs.Product) {
    resp := ResponseStruct{"bad product"}
    js, err := json.Marshal(resp)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    w.WriteHeader(http.StatusBadRequest)
    w.Write(js)
    return
  }

  if !checkPlan(rqs.Plan) {
    resp := ResponseStruct{"bad plan"}
    js, err := json.Marshal(resp)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    w.WriteHeader(http.StatusBadRequest)
    w.Write(js)
    return
  }

  if !checkRegion(rqs.Region) {
    resp := ResponseStruct{"bad region"}
    js, err := json.Marshal(resp)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    w.WriteHeader(http.StatusBadRequest)
    w.Write(js)
    return
  }

  if resourceExistsAndResponseCreated(r.Method, rqs, w) {
    return
  }

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
  return
}

func credentialsHandler(w http.ResponseWriter, r *http.Request) {
  return
}

func ssoHandler(w http.ResponseWriter, r *http.Request) {
  return
}

func seed() string {
	rand.Seed(time.Now().UTC().UnixNano())
	result := rand.Intn(100)
  return fmt.Sprintf("%d", result)
}

func checkPlan(plan string) bool {
  plans := [2]string{"small", "large"}

  for i := range plans {
    if plan == plans[i] {
      return true
    }
  }

  return false
}

func checkProduct(product string) bool {
  products := [1]string{"bonnets"}

  for i := range products {
    if product == products[i] {
      return true
    }
  }

  return false
}

func checkRegion(region string) bool {
  regions  := [1]string{"aws::us-east-1"}

  for i := range regions {
    if region == regions[i] {
      return true
    }
  }

  return false
}

func resourceExistsAndResponseCreated(method string, rqs RequestStruct, w http.ResponseWriter) bool {
  existingData, dataRetrieved := db[rqs.Id]
  newData := string(convertRequestToJson(rqs))
  newDataMatchesOldData := existingData == newData

  if dataRetrieved && (method == "POST" || method == "PUT") {
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
