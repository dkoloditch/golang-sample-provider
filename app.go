package main

import (
  "fmt"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "golang.org/x/oauth2"
  "os"
  "encoding/json"
  "path"
)

var (
  // simple in-memory db
  db = Database{
    Resources: make(map[string]string),
    Credentials: make(map[string]string),
  }
  MASTER_KEY = os.Getenv("MASTER_KEY")
  CLIENT_ID = os.Getenv("CLIENT_ID")
  CLIENT_SECRET = os.Getenv("CLIENT_SECRET")
  CONNECTOR_URL = os.Getenv("CONNECTOR_URL")
  oac = &oauth2.Config{
    ClientID: CLIENT_ID,
    ClientSecret: CLIENT_SECRET,
    Scopes:       []string{},
    Endpoint: oauth2.Endpoint{
        AuthURL:  CONNECTOR_URL,
        TokenURL: CONNECTOR_URL + "/oauth/tokens",
    },
  }
)

func main() {
  router := mux.NewRouter().StrictSlash(true)

  router.HandleFunc("/dashboard", dashboardHandler).Methods("GET")

  router.HandleFunc("/v1/resources/{id}", createResourcessHandler).Methods("PUT")
  router.HandleFunc("/v1/resources/{id}", updateResourcessHandler).Methods("PATCH")
  router.HandleFunc("/v1/resources/{id}", deleteResourcessHandler).Methods("DELETE")

  router.HandleFunc("/v1/credentials/{id}", createCredentialsHandler).Methods("PUT")
  router.HandleFunc("/v1/credentials/{id}", deleteCredentialsHandler).Methods("DELETE")

  router.HandleFunc("/v1/sso", ssoHandler).Methods("GET")

  log.Fatal(http.ListenAndServe(":4567", router))
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
  return
}

func createResourcessHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  bodyBuffer, rqs := GetBodyBufferAndResources(r)
  id := rqs.Id

  if SignatureIsNotValidAndResponseCreated(r, w, bodyBuffer) { return }

  if ProductIsNotValidAndResponseCreated(rqs.Product, w) { return }

  if PlanIsNotValidAndResponseCreated(rqs.Plan, w) { return }

  if RegionIsNotValidAndResponseCreated(rqs.Region, w) { return }

  if ResourceAlreadyExistsAndResponseCreated(rqs, w, id) { return }

  if ValidCreateRequestAndResponseCreated(rqs, w) { return }
}

func updateResourcessHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  // since the id is only passed via URL with PATCH requests, we set this here
  // and provide it to the relevant methods below.
  _, id := path.Split(r.URL.Path)
  bodyBuffer, rqs := GetBodyBufferAndResources(r)

  if SignatureIsNotValidAndResponseCreated(r, w, bodyBuffer) { return }

  if PlanIsNotValidAndResponseCreated(rqs.Plan, w) { return }

  if ResourceDoesNotExistAndResponseCreated(rqs, w, id) { return }

  if ValidUpdateRequestAndResponseCreated(rqs, w, id) { return }
}

func deleteResourcessHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  _, id := path.Split(r.URL.Path)
  bodyBuffer, rqs := GetBodyBufferAndResources(r)

  if SignatureIsNotValidAndResponseCreated(r, w, bodyBuffer) { return }

  if ResourceDoesNotExistAndResponseCreated(rqs, w, id) { return }

  if ResourceDeletedAndResponseCreated(w, id) { return }

  return
}

func createCredentialsHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  bodyBuffer, rqs := GetBodyBufferAndCredentials(r)

  if SignatureIsNotValidAndResponseCreated(r, w, bodyBuffer) { return }

  if InvalidResourceIdAndResponseCreated(w, rqs) { return }

  if ProvisionCredentialsAndResponseCreated(w, rqs) { return }

  return
}

func deleteCredentialsHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  _, id := path.Split(r.URL.Path)

  if CredentialsDoNotExistAndResponseCreated(w, id) { return }

  if CredentialsDeletedAndResponseCreated(w, id) { return }
}

func ssoHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  code := r.URL.Query().Get("code")
  // resource_id := r.URL.Query().Get("resource_id")
  // url := oac.AuthCodeURL(CONNECTOR_URL)
  _, err := oac.Exchange(oauth2.NoContext, code) // token, err
  errReformatted := fmt.Errorf("%v", err) // avoids blowup

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
