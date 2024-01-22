package api

import (
	"backend/db"
	"backend/types"
	"backend/util"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// client must be authenticated to create a furniture listing
	sessionManager := GetSessionManager()
	_, loggedIn := sessionManager.IsLoggedIn(r)
	if !loggedIn {
		http.Error(w, "You must be logged in to create a furniture listing", http.StatusUnauthorized)
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
This endpoint handler returns the furniture listing given the listingID
is provided in the request URL

200 - furniture listing found
500 - furniture listing not found in DB
*/
func (s *Server) HandleGetFurniture(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Request must be a GET request", http.StatusMethodNotAllowed)
		return
	}

	// get id from url query params
	id := r.URL.Query().Get("listingid")

	_, err := db.FindByIDInListingsCollection(id)
	if err != nil {
		http.Error(
			w,
			"Furniture listing with provided listingID not found",
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))

}

/*
This endpoint handler gets all the listed furnitures and returns
it back to client in response

- Note: Wouldn't it be bad to return EVERY single document in the listings
collection at the same time for each request?

For now, this endpoint returns every single document in the listings collection
at the same time for each request
*/
func (s *Server) HandleGetFurnitures(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Request must be a GET request", http.StatusMethodNotAllowed)
		return
	}

	collection := db.GetCollection("listings")
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		http.Error(w, "Error getting listings", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var listings []types.FurnitureListing
	for cursor.Next(context.Background()) {
		var listing types.FurnitureListing
		if err := cursor.Decode(&listing); err != nil {
			log.Printf("Error decoding document: %v", err)
			continue
		}
		listings = append(listings, listing)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Error iterating over furniture listings", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(listings)
	if err != nil {
		http.Error(w, "Error encoding response data into JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}
