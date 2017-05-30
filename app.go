package main

import (
  "fmt"
  // "html"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  // "golang.org/x/oauth2"
  // "os"
  "math/rand"
  "time"
  "encoding/json"
)

func main() {
  // MASTER_KEY := os.Getenv("MASTER_KEY")
  // CLIENT_ID := os.Getenv("CLIENT_ID")
  // CLIENT_SECRET := os.Getenv("CLIENT_SECRET")
  // CONNECTOR_URL := os.Getenv("CONNECTOR_URL")

  // fmt.Println(MASTER_KEY)
  // fmt.Println(CLIENT_ID)
  // fmt.Println(CLIENT_SECRET)
  // fmt.Println(CONNECTOR_URL)

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
  fmt.Println("*** dashboardHandler ***")
  return
}

func resourcesHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("*** resourcesHandler ***")

  // decode request body
  decoder := json.NewDecoder(r.Body)
  var rqs RequestStruct
  decoder.Decode(&rqs)
  fmt.Println(rqs)
  fmt.Println("\n")

  if !checkProduct(rqs.Product) {
    resp := ResponseStruct{"bad product"}
    js, err := json.Marshal(resp)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    w.Header().Set("Content-Type", "application/json")
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

    w.Header().Set("Content-Type", "application/json")
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

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    w.Write(js)
    return
  }

  // get the random number and create json response
  result := seed()
  resp := ResponseStruct{fmt.Sprintf("%d", result)}
  js, err := json.Marshal(resp)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  w.Write(js)
  return
}

func credentialsHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("*** credentialsHandler ***")
  return
}

func ssoHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("*** ssoHandler ***")
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

type ResponseStruct struct {
  Message string `json:"message"`
}

type RequestStruct struct {
  Id string
  Product string
  Plan string
  Region string
}
