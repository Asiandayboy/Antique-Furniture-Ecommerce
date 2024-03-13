package api

import (
	"backend/db"
	"backend/types"
	"backend/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ErrAddressNoChanges = "Empty fields; no address changes provided"
)

/*
Represent changes a user makes to their account and is also used
as the response to the client to inform them of their changes
*/
type AccountEdit struct {
	NewPassword string `bson:"password,omitempty" json:"newPassword,omitempty"`
	NewPhone    string `bson:"phone,omitempty" json:"newPhone,omitempty"`
	NewEmail    string `bson:"email,omitempty" json:"newEmail,omitempty"`
}

/*
Represents user changes to an existing address
*/
type AddressUpdate struct {
	NewState   string `bson:"state,omitempty" json:"state"`
	NewCity    string `bson:"city,omitempty" json:"city"`
	NewStreet  string `bson:"street,omitempty" json:"street"`
	NewZipCode string `bson:"zipCode,omitempty" json:"zipCode"`
	NewDefault bool   `bson:"default,omitempty" json:"default"`
}

func (a AddressUpdate) IsEmpty() bool {
	// ignore NewDefault field because its zero value is false by default
	if a.NewCity == "" &&
		a.NewState == "" &&
		a.NewStreet == "" &&
		a.NewZipCode == "" {
		return true
	}

	return false
}

/*
Represents the JSON input format for PUT /account/address
*/
type AddressUpdateInput struct {
	AddressID string        `json:"addressID"` // The hex string ID of the document to edit
	Changes   AddressUpdate `json:"changes"`   // User changes for the the address
}

/*
Retrieve account information
*/
func (s *Server) HandleAccountGET(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(SessionKey).(*Session)

	var userInfo types.User
	usersCollection := db.GetCollection("users")
	err := usersCollection.FindOne(context.Background(), bson.M{
		"_id": session.Store["userid"],
	}).Decode(&userInfo)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to fetch account: %s\n", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	// prepare data to send back to client
	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

/*
Update current account information
*/
func (s *Server) HandleAccountPUT(w http.ResponseWriter, r *http.Request) {
	var changes AccountEdit
	if err := util.ReadJSONReq[AccountEdit](r, &changes); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	/*
		If the user is changing their password, the input will contain the plain
		text password. We must hash it before saving it
	*/
	if changes.NewPassword != "" {
		// hash new password
		hashedPass, err := util.HashPassword(changes.NewPassword)
		if err != nil {
			http.Error(w, "Failed to hash new password", http.StatusInternalServerError)
			return
		}

		changes.NewPassword = hashedPass
	}

	session := r.Context().Value(SessionKey).(*Session)

	usersCollection := db.GetCollection("users")
	_, err := usersCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": session.Store["userid"]},
		bson.M{"$set": changes},
	)

	if err != nil {
		http.Error(w, "Failed to update account information", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

// Used to retrieve all of the user's addresses
func (s *Server) HandleAddressGET(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(SessionKey).(*Session)

	shippingAddrColl := db.GetCollection("shippingAddresses")
	cursor, err := shippingAddrColl.Find(
		context.Background(),
		bson.M{"userid": session.Store["userid"]},
	)
	if err != nil {
		http.Error(w, "Failed to fetch shipping addresses", http.StatusInternalServerError)
		return
	}

	var documents []types.ShippingAddress
	err = cursor.All(context.Background(), &documents)
	if err != nil {
		http.Error(w, "Failed to cursor.All", http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(documents)
	if err != nil {
		http.Error(w, "Failed to encode into JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// Used to create a new address
func (s *Server) HandleAddressPOST(w http.ResponseWriter, r *http.Request) {
	var address types.ShippingAddress
	err := util.ReadJSONReq[types.ShippingAddress](r, &address)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	session := r.Context().Value(SessionKey).(*Session)

	// set userID before updatting users document
	address.UserID = session.Store["userid"].(primitive.ObjectID)

	addressesCollection := db.GetCollection("shippingAddresses")
	_, err = addressesCollection.InsertOne(context.Background(), address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

// Used to edit and update an existing address
func (s *Server) HandleAddressPUT(w http.ResponseWriter, r *http.Request) {
	var changes AddressUpdateInput
	err := util.ReadJSONReq[AddressUpdateInput](r, &changes)
	if err != nil {
		http.Error(w, "Failed decode from JSON", http.StatusBadRequest)
		return
	}

	// validate input to ensure fields are not empty
	if changes.Changes.IsEmpty() {
		http.Error(w, ErrAddressNoChanges, http.StatusBadRequest)
		return
	}

	addressID, err := primitive.ObjectIDFromHex(changes.AddressID)
	if err != nil {
		http.Error(w, primitive.ErrInvalidHex.Error(), http.StatusBadRequest)
		return
	}

	addressesCollection := db.GetCollection("shippingAddresses")
	res, err := addressesCollection.UpdateByID(
		context.Background(),
		addressID,
		bson.M{"$set": changes.Changes},
	)
	if err != nil {
		http.Error(w, "Failed to update record", http.StatusInternalServerError)
		return
	}
	if res.MatchedCount == 0 {
		http.Error(
			w,
			"Document with provided addressID does not exist",
			http.StatusBadRequest,
		)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

// Used to delete an address
func (s *Server) HandleAddressDELETE(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(SessionKey).(*Session)

	var addressID string = r.PathValue("addressID")

	objID, err := primitive.ObjectIDFromHex(addressID)
	if err != nil {
		http.Error(w, primitive.ErrInvalidHex.Error(), http.StatusBadRequest)
		return
	}

	shippingAddrColl := db.GetCollection("shippingAddresses")
	res, err := shippingAddrColl.DeleteOne(
		context.Background(),
		bson.M{"_id": objID, "userid": session.Store["userid"]},
	)
	if err != nil {
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}
	if res.DeletedCount == 0 {
		http.Error(w, "Could not find document with ID", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))

}

// Fetches the user's entire purchase history
func (s *Server) HandlePurchaseHistory(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(SessionKey).(*Session)

	receiptsCollection := db.GetCollection("receipts")
	cursor, err := receiptsCollection.Find(
		context.Background(),
		bson.M{"userid": session.Store["userid"]},
	)
	if err != nil {
		http.Error(w, "Failed to fetch purchase history", http.StatusInternalServerError)
		return
	}

	var receipts []Receipt
	err = cursor.All(context.Background(), &receipts)
	if err != nil {
		http.Error(w, "Failed to cursor.All", http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(receipts)
	if err != nil {
		http.Error(w, "Failed to encode into JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

/*
Returns a specified purchase history item to the client
*/
func (s *Server) HandlePurchaseHistoryItem(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(SessionKey).(*Session)

	var id string = r.PathValue("orderID")
	orderID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, primitive.ErrInvalidHex.Error(), http.StatusBadRequest)
		return
	}

	receipsCollection := db.GetCollection("receipts")

	// find specified order
	var order Receipt
	res := receipsCollection.FindOne(
		context.Background(),
		bson.M{"_id": orderID, "userid": session.Store["userid"]},
	).Decode(&order)
	if res != nil {
		http.Error(w, res.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Failed to encode into JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
}
