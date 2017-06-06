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

	SetHeadersAndResponse("bad product", http.StatusBadRequest, RESPONSE_TYPE_GENERAL, w) // 400

	return true
}

func PlanIsNotValid(plan string, w http.ResponseWriter) bool {
	for i := range plans {
		if plan == plans[i] {
			return false
		}
	}

	SetHeadersAndResponse("bad plan", http.StatusBadRequest, RESPONSE_TYPE_GENERAL, w) // 400

	return true
}

func RegionIsNotValid(region string, w http.ResponseWriter) bool {
	for i := range regions {
		if region == regions[i] {
			return false
		}
	}

	SetHeadersAndResponse("bad region", http.StatusBadRequest, RESPONSE_TYPE_GENERAL, w) // 400

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
			message := ""

			SetHeadersAndResponse(message, http.StatusNoContent, RESPONSE_TYPE_GENERAL, w)

			return true
		} else {
			// different content results in conflict
			message := "resource already exists"

			SetHeadersAndResponse(message, http.StatusConflict, RESPONSE_TYPE_GENERAL, w)

			return true
		}
	}

	return false
}

func ResourceDoesNotExist(w http.ResponseWriter, id string) bool {
	_, dataRetrieved := db.Resources[id]

	if !dataRetrieved {
		// non existing resource
		message := "no such resource"

		SetHeadersAndResponse(message, http.StatusNotFound, RESPONSE_TYPE_GENERAL, w)

		return true
	}

	return false
}

func ResourceCreated(rqs Resources, w http.ResponseWriter) bool {
	// get the random number and create json response
	result := Seed()
	message := fmt.Sprintf("%d", result)
	rqsData := Resources{rqs.Id, rqs.Product, rqs.Plan, rqs.Region, result}
	data, err := json.Marshal(rqsData)

	IssueResponseIfErrorOccurs(err, w)

	// add to db
	db.Resources[rqs.Id] = string(data)

	SetHeadersAndResponse(message, http.StatusCreated, RESPONSE_TYPE_GENERAL, w)

	return true
}

func ResourceUpdated(rqs Resources, w http.ResponseWriter, id string) bool {
	// get the random number and create json response
	result := Seed()
	message := fmt.Sprintf("%d", result)
	rqsData := Resources{rqs.Id, rqs.Product, rqs.Plan, rqs.Region, result}
	data, err := json.Marshal(rqsData)

	IssueResponseIfErrorOccurs(err, w)

	// remove previous data and add new data to db
	// @TODO: this is obviously a hack since we're completely replacing the data
	// (PUT) rather than modifying it (PATCH). fix this eventually.
	delete(db.Resources, rqs.Id)
	db.Resources[rqs.Id] = string(data)

	SetHeadersAndResponse(message, http.StatusOK, RESPONSE_TYPE_GENERAL, w)

	return false
}

func ResourceDeleted(w http.ResponseWriter, id string) bool {
	delete(db.Resources, id)

	message := ""

	SetHeadersAndResponse(message, http.StatusNoContent, RESPONSE_TYPE_GENERAL, w)

	return true
}

func InvalidResourceForCredential(w http.ResponseWriter, rqs CredentialsRequest) bool {
	resource := db.Resources[rqs.ResourceId]

	if resource == "" {
		message := "no such resource"

		SetHeadersAndResponse(message, http.StatusNotFound, RESPONSE_TYPE_GENERAL, w)

		return true
	}

	return false
}
