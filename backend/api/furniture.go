package api

import (
	"backend/db"
	"backend/types"
	"backend/util"
	"net/http"

	"github.com/google/uuid"
)

type listFormError struct {
	message string
}

func (l listFormError) Error() string {
	return l.message
}

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
*/
func (s *Server) HandleListFurniture(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Request must be a POST request", http.StatusBadRequest)
		return
	}

	// decode request body into struct
	var newListing types.FurnitureListing
	err := util.ReadJSONReq[types.FurnitureListing](r, &newListing)
	if err != nil {
		http.Error(w, "Failed to decode into JSON", http.StatusInternalServerError)
		return
	}

	// validate form inputs
	validated, errMsg := ValidateListFormFields(newListing)
	if !validated {
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// apply a unique ID to the listing
	listingId, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, "Failed to apply ID to new listing", http.StatusInternalServerError)
		return
	}
	newListing.Id = listingId.String()

	// save new listing in database
	result, _ := db.InsertIntoListingsCollection(newListing)

	// return success code, "success" msg

	http.Error(w, "something", http.StatusInternalServerError)
}

/*
This endpoint handler gets all the listed furnitures and returns
it back to client in response
*/
func HandleGetFurnitures(w http.ResponseWriter, r *http.Request) {

}
