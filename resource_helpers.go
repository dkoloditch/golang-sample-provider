package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func ProductIsNotValid(product string, w http.ResponseWriter) bool {
	for i := range products {
		if product == products[i] {
			return false
		}
	}

	HandleResponse("bad product", http.StatusBadRequest, w) // 400

	return true
}

func PlanIsNotValid(plan string, w http.ResponseWriter) bool {
	for i := range plans {
		if plan == plans[i] {
			return false
		}
	}

	HandleResponse("bad plan", http.StatusBadRequest, w) // 400

	return true
}

func RegionIsNotValid(region string, w http.ResponseWriter) bool {
	for i := range regions {
		if region == regions[i] {
			return false
		}
	}

	HandleResponse("bad region", http.StatusBadRequest, w) // 400

	return true
}

func NoDifferenceInContent(rqs1 Resources, rqs2 Resources) bool {
	productMatch := rqs1.Product == rqs2.Product
	planMatch := rqs1.Plan == rqs2.Plan
	regionMatch := rqs1.Region == rqs2.Region

	return productMatch && planMatch && regionMatch
}

func ResourceAlreadyExists(rqs Resources, w http.ResponseWriter, id string) bool {
	// this function is only used for create / POST attempts

	existingDataRetrieved, dataRetrieved := db.Resources[id]
	existingDataBytes := []byte(existingDataRetrieved)
	existingDataBuffer := bytes.NewBuffer(existingDataBytes)
	existingDataRqsStruct := GetResources(existingDataBuffer)
	resourceAlreadyExists := existingDataRqsStruct.Id == id
	NoDifferenceInContent := NoDifferenceInContent(existingDataRqsStruct, rqs)

	// @TODO this can probably be refactored to not account for Method given
	// appropriate routes
	if dataRetrieved {
		// same content acts as created
		if resourceAlreadyExists && NoDifferenceInContent {
			// @TODO: respond with appropriate random number?
			resp := ResponseStruct{""}
			js, err := json.Marshal(resp)

			IssueResponseIfErrorOccurs(err, w)

			w.WriteHeader(http.StatusNoContent) // 204
			w.Write(js)

			return true
		} else {
			// different content results in conflict
			resp := ResponseStruct{"resource already exists"}
			js, err := json.Marshal(resp)

			IssueResponseIfErrorOccurs(err, w)

			w.WriteHeader(http.StatusConflict) // 409
			w.Write(js)

			return true
		}
	}

	return false
}

func ResourceDoesNotExist(w http.ResponseWriter, id string) bool {
	_, dataRetrieved := db.Resources[id]

	if !dataRetrieved {
		// non existing resource
		resp := ResponseStruct{"no such resource"}
		js, err := json.Marshal(resp)

		IssueResponseIfErrorOccurs(err, w)

		w.WriteHeader(http.StatusNotFound) // 404
		w.Write(js)

		return true
	}

	return false
}

func ResourceCreated(rqs Resources, w http.ResponseWriter) bool {
	// @TODO: can this be abstracted further?
	// @TODO: save random number in db with data

	// get the random number and create json response
	result := Seed()
	resp := ResponseStruct{fmt.Sprintf("%d", result)}
	rqsData := Resources{rqs.Id, rqs.Product, rqs.Plan, rqs.Region, result}
	js, err := json.Marshal(resp)
	data, err := json.Marshal(rqsData)

	IssueResponseIfErrorOccurs(err, w)

	// add to db
	db.Resources[rqs.Id] = string(data)

	w.WriteHeader(http.StatusCreated) // 201
	w.Write(js)

	return true
}

func ResourceUpdated(rqs Resources, w http.ResponseWriter, id string) bool {
	// get the random number and create json response
	result := Seed()
	resp := ResponseStruct{fmt.Sprintf("%d", result)}
	rqsData := Resources{rqs.Id, rqs.Product, rqs.Plan, rqs.Region, result}
	js, err := json.Marshal(resp)
	data, err := json.Marshal(rqsData)

	IssueResponseIfErrorOccurs(err, w)

	// remove previous data and add new data to db
	// @TODO: this is obviously a hack since we're completely replacing the data
	// (PUT) rather than modifying it (PATCH). fix this eventually.
	delete(db.Resources, rqs.Id)
	db.Resources[rqs.Id] = string(data)

	w.WriteHeader(http.StatusOK) // 200
	w.Write(js)

	return false
}

func ResourceDeleted(w http.ResponseWriter, id string) bool {
	delete(db.Resources, id)

	resp := &ResponseStruct{}
	js, err := json.Marshal(resp)

	IssueResponseIfErrorOccurs(err, w)

	w.WriteHeader(http.StatusNoContent) // 204
	w.Write(js)

	return true
}

func InvalidResourceForCredential(w http.ResponseWriter, rqs CredentialsRequest) bool {
	resource := db.Resources[rqs.ResourceId]

	if resource == "" {
		resp := CredentialsResponse{
			Message: "no such resource",
		}
		js, err := json.Marshal(resp)

		IssueResponseIfErrorOccurs(err, w)

		w.WriteHeader(http.StatusNotFound)
		w.Write(js)

		return true
	}

	return false
}
