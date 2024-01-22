package api

import (
	"backend/db"
	"backend/types"
	"backend/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// this has to be the ugliest thiing ever
func ValidateListFormFields(listing types.FurnitureListing) (bool, string) {
	if listing.Title == "" {
		return false, "Title not provided"
	} else if listing.Condition == "" {
		return false, "Condition not provided"
	} else if listing.Cost == 0 {
		return false, "Cost not provided"
	} else if listing.Description == "" {
		return false, "Description not provided"
	} else if listing.Type == "" {
		return false, "Furniture type not provided"
	} else if listing.Material == "" {
		return false, "Material not provided"
	} else if listing.Style == "" {
		return false, "Style not provided"
	} else if len(listing.Images) == 0 {
		return false, "Images not provided"
	}
	return true, ""
}

/*
This endpoint processes a new furniture listing request

Successful requests return 200 status code and the new listingID
*/
func (s *Server) HandleListFurniture(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Request must be a POST request", http.StatusMethodNotAllowed)
		return
	}

	// decode request body into struct
	var newListing types.FurnitureListing
	err := util.ReadJSONReq[types.FurnitureListing](r, &newListing)
	if err != nil {
		http.Error(w, "Failed to decode JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// validate form inputs
	validated, errMsg := ValidateListFormFields(newListing)
	if !validated {
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// save new listing in database
	result, err := db.InsertIntoListingsCollection(newListing)
	if err != nil {
		http.Error(w, "Failed to insert listing into database", http.StatusConflict)
		return
	}

	// insertedID is of type primitive.ObjectID, which is type [12]byte
	insertedId := result.InsertedID.(primitive.ObjectID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(insertedId.Hex()))

}

/*
This endpoint handler gets all the listed furnitures and returns
it back to client in response
*/
func HandleGetFurnitures(w http.ResponseWriter, r *http.Request) {

}
