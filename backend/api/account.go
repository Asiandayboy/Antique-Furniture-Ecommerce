package api

import (
	"backend/db"
	"backend/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

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

}

/*
This handler processes requests for fetching the client's addresses,
creating addresses, updating them, or deleteing them
*/
func (s *Server) HandleAddresses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleAddressesGet(w, r)
	case "POST":
		handleAddressesPost(w, r)
	case "PUT":
		handleAddressesPut(w, r)
	case "DELETE":
		handleAddressesDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Used to retrieve addresses
func handleAddressesGet(w http.ResponseWriter, r *http.Request) {

}

// Used to create a new address
func handleAddressesPost(w http.ResponseWriter, r *http.Request) {

}

// Used to edit and update an existing address
func handleAddressesPut(w http.ResponseWriter, r *http.Request) {

}

// Used to delete an address
func handleAddressesDelete(w http.ResponseWriter, r *http.Request) {

}
