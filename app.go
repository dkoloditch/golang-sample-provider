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
  return
}

func resourcesHandler(w http.ResponseWriter, r *http.Request) {
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

type ResponseStruct struct {
  Message string
}

type RequestStruct struct {
  Id string
  Product string
  Plan string
  Region string
}
