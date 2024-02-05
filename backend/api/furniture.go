package api

import (
	"backend/db"
	"backend/types"
	// "backend/util"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ErrListFormNoCondition       = "Furniture condition not provided"
	ErrListFormNoCost            = "Furniture cost not provided"
	ErrListFormNoDescription     = "Furniture description not provided"
	ErrListFormNoImages          = "Furniture images not provided"
	ErrListFormNoMaterial        = "Furniture material not provided"
	ErrListFormNoStyle           = "Furniture style not provided"
	ErrListFormNoTitle           = "Furniture title not provided"
	ErrListFormNoType            = "Furniture type not provided"
	ErrListFormEveryFieldMissing = "Every field is missing"
)

const NUMBER_OF_LIST_FORM_FIELDS = 8 // removed images
// const NUMBER_OF_LIST_FORM_FIELDS = 7

type ListFormErrors struct {
	FormErrors []string `json:"formErrors"`
	length     int8
}

func (l ListFormErrors) Error() string {
	if l.length == NUMBER_OF_LIST_FORM_FIELDS {
		return ErrListFormEveryFieldMissing
	}
	jsonData, _ := json.Marshal(l.FormErrors)
	return string(jsonData)
}

/*
Returns nil or an error if any of the fields for the
FurnitureListing form is empty or not provided

Calling .Error() on the error will return a string of the array
of each error in alphabetical order, where <ErrListFormNoCondition>
is the first, <ErrListFormNoCost> second, and so on.

If every field is missing, then .Error() will return <ErrListFormEveryFieldMissing>
*/
func ValidateListFormFields(listing types.FurnitureListing) error {
	formErrs := make([]string, 0, 8)
	var length int8 = 0
	if listing.Condition == "" {
		formErrs = append(formErrs, ErrListFormNoCondition)
		length++
	}
	if listing.Cost == 0 {
		formErrs = append(formErrs, ErrListFormNoCost)
		length++
	}
	if listing.Description == "" {
		formErrs = append(formErrs, ErrListFormNoDescription)
		length++
	}
	if len(listing.Images) == 0 {
		formErrs = append(formErrs, ErrListFormNoImages)
		length++
	}
	if listing.Material == "" {
		formErrs = append(formErrs, ErrListFormNoMaterial)
		length++
	}
	if listing.Style == "" {
		formErrs = append(formErrs, ErrListFormNoStyle)
		length++
	}
	if listing.Title == "" {
		formErrs = append(formErrs, ErrListFormNoTitle)
		length++
	}
	if listing.Type == "" {
		formErrs = append(formErrs, ErrListFormNoType)
		length++
	}

	if length == 0 {
		return nil
	}
	return ListFormErrors{
		FormErrors: formErrs,
		length:     length,
	}
}

/*
This endpoint processes a new furniture listing request

If an error occurs during the processing, an error message string
will be returned in the response.

If no errors occur, the ID hex string of the ObjectID of the new
furniture listing document will be returned in the response.

TODO: MUST CHANGE TO ACCEPT FILE UPLOAD FOR IMAGES (MULTIPART) 2/3
*/
func (s *Server) HandleListFurniture(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error()+" --> err1", http.StatusBadRequest)
		return
	}

	// decode request body into struct
	jsonData := r.MultipartForm.Value["json_data"]
	if len(jsonData) == 0 {
		http.Error(w, "JSON data not provided in the form", http.StatusBadRequest)
		return
	}

	var newListing types.FurnitureListing
	if err := json.Unmarshal([]byte(jsonData[0]), &newListing); err != nil {
		http.Error(w, "Failed to unmarshal JSON data", http.StatusBadRequest)
		return
	}

	// get images
	files := r.MultipartForm.File["furniture_images"]
	for _, file := range files {
		// open file
		fileReader, err := file.Open()
		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer fileReader.Close()

		// read file
		fileData, err := io.ReadAll(fileReader)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		// add image data to newListing
		newListing.Images = append(newListing.Images, fileData)
	}

	// validate form inputs
	err = ValidateListFormFields(newListing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// add userID to newListing
	session := r.Context().Value(SessionKey).(*Session)
	newListing.UserID = session.Store["userid"].(primitive.ObjectID)

	// save new listing in database
	listingsCollection := db.GetCollection("listings")
	result, err := listingsCollection.InsertOne(context.Background(), newListing)
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
500 - furniture listing not found in DB or 400??
*/
func (s *Server) HandleGetFurniture(w http.ResponseWriter, r *http.Request) {
	// get id from url query params
	// listingid param might not be set; check for that 1/26
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
