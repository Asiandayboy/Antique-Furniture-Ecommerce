package api

import (
	"backend/db"
	"backend/types"
	"backend/util"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
This type is used to represent changes a user makes to their account and
is used as the response to the client to inform them of the changes
*/
type AccountEdit struct {
	NewPassword string `bson:"password,omitempty" json:"newPassword,omitempty"`
	NewPhone    string `bson:"phone,omitempty" json:"newPhone,omitempty"`
	NewEmail    string `bson:"email,omitempty" json:"newEmail,omitempty"`
}

/*
Retrieve account information
*/
func (s *Server) HandleAccountGET(w http.ResponseWriter, r *http.Request) {
	log.Println("\x1b[35mENDPOINT HIT -> /account GET\x1b[0m")

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
	log.Println("\x1b[35mENDPOINT HIT -> /account PUT\x1b[0m")

	var changes AccountEdit
	if err := util.ReadJSONReq[AccountEdit](r, &changes); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	/*
		If the user is changing their password, the request will contain the plain
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

// Used to retrieve addresses
func (s *Server) HandleAddressGET(w http.ResponseWriter, r *http.Request) {
	log.Println("\x1b[35mENDPOINT HIT -> /account/address GET\x1b[0m")

	http.Error(w, "test", http.StatusInternalServerError)
}

// Used to create a new address
func (s *Server) HandleAddressPOST(w http.ResponseWriter, r *http.Request) {
	log.Println("\x1b[35mENDPOINT HIT -> /account/address POST\x1b[0m")

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
	log.Println("\x1b[35mENDPOINT HIT -> /account/address PUT\x1b[0m")
	http.Error(w, "test", http.StatusInternalServerError)

}

// Used to delete an address
func (s *Server) HandleAddressDELETE(w http.ResponseWriter, r *http.Request) {
	log.Println("\x1b[35mENDPOINT HIT -> /account/address DELETE\x1b[0m")
	http.Error(w, "test", http.StatusInternalServerError)

}
