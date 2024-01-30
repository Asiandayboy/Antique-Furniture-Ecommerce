package tests

import (
	"backend/api"
	"backend/db"
	"backend/types"
	"backend/util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func trimSpaceAndNewline(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "\n", ""))
}

func TestHandleSignup(t *testing.T) {
	db.Init()
	defer db.Close()
	server := api.NewServer(":3000")
	server.Post("/signup", server.HandleSignup)

	tests := []struct {
		name               string
		method             string
		payload            string
		expectedResMsg     string
		expectedStatusCode int
	}{
		// { // test valid signup; if this test fails, it's bc the payload is already in the db; just delete it and test again
		// 	name:               "Test 1",
		// 	method:             "POST",
		// 	payload:            `{"username": "testuser1", "password": "testpassword1", "email": "test@gmail.com"}`,
		// 	expectedResMsg:     "success",
		// 	expectedStatusCode: http.StatusOK,
		// },
		{ // test invalid method
			name:               "Test 2",
			method:             "GET",
			payload:            `{"username": "testuser1", "password": "testpassword1", "email": "test@gmail.com"}`,
			expectedResMsg:     "Request must be a POST request",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{ // testing invalid json decode
			name:               "Test 3",
			method:             "POST",
			payload:            `{"username": "testuser1", "password": "testpassword1, "email": "test@gmail.com"}`,
			expectedResMsg:     "Could not decode request body into JSON",
			expectedStatusCode: http.StatusBadRequest,
		},
		{ // testing username against existing username -> "bob"
			name:               "Test 4",
			method:             "POST",
			payload:            `{"username": "bob", "password": "testpassword1", "email": "test@gmail.com"}`,
			expectedResMsg:     "Username is taken",
			expectedStatusCode: http.StatusConflict,
		},
		{ // testing email against existing email -> "johnsmith@gmail.com"
			name:               "Test 5",
			method:             "POST",
			payload:            `{"username": "testuser123", "password": "testpassword1", "email": "johnsmith@gmail.com"}`,
			expectedResMsg:     "Email is taken",
			expectedStatusCode: http.StatusConflict,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/signup", bytes.NewBufferString(tc.payload))

			// response recorder response the response
			w := httptest.NewRecorder()

			// hit endpoint
			server.Mux.ServeHTTP(w, r)

			if status := w.Code; status != tc.expectedStatusCode {
				t.Errorf("Expected code: %v, got: %v", tc.expectedStatusCode, status)
			}
			if res := w.Body; trimSpaceAndNewline(res.String()) != tc.expectedResMsg {
				t.Errorf("Expected ResMsg: %v, got: %v", tc.expectedResMsg, res.String())
			}

		})
	}
}

/*
When running this test, make sure to run
TestHandleSignup so that the mock data gets added.
This test uses the the mock data from that test

If you don't this test will not work properly
*/
func TestHandleLogin(t *testing.T) {
	db.Init()
	defer db.Close()
	server := api.NewServer(":3000")
	server.Post("/login", server.HandleLogin)

	tests := []struct {
		name               string
		method             string
		payload            string
		expectedResMsg     string
		expectedStatusCode int
	}{
		{ // test valid login (TestHandleSignup should be ran first to test this first case)
			name:               "Test 1",
			method:             "POST",
			payload:            `{"username": "testuser1", "password": "testpassword1"}`,
			expectedResMsg:     "success",
			expectedStatusCode: http.StatusOK,
		},
		{ // test invalid method
			name:               "Test 2",
			method:             "GET",
			payload:            `{"username": "bob", "password": "bob"}`,
			expectedResMsg:     "Request must be a POST request",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{ // test invalid login/invalid username
			name:               "Test 3",
			method:             "POST",
			payload:            `{"username": "testuser123", "password": "testpassword1"}`,
			expectedResMsg:     "Invalid login",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{ // test invalid login/incorrect password
			name:               "Test 4",
			method:             "POST",
			payload:            `{"username": "testuser1", "password": "testpassword123"}`,
			expectedResMsg:     "Invalid login",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{ // test invalid json decode
			name:               "Test 5",
			method:             "POST",
			payload:            `{"username": "bob", "password: "bob123"}`,
			expectedResMsg:     "Could not decode request body into JSON",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/login", bytes.NewBufferString(tc.payload))

			w := httptest.NewRecorder()
			server.Mux.ServeHTTP(w, r)

			if status := w.Code; status != tc.expectedStatusCode {
				t.Errorf("Expected code: %v, got: %v", tc.expectedStatusCode, status)
			}
			if res := w.Body; trimSpaceAndNewline(res.String()) != tc.expectedResMsg {
				t.Errorf("Expected ResMsg: %v, got: %v", tc.expectedResMsg, res.String())
			}
		})
	}
}

func TestHandleLogout(t *testing.T) {
	sessionManager := api.GetSessionManager()

	// creating fake loggedIn client #1
	session1, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "test1",
	})
	if err != nil {
		t.Fatal("Failed to create fake session1")
		return
	}

	// creating fake loggedIn client #2
	session2, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "test2",
	})
	if err != nil {
		t.Fatal("Failed to create fake sesion2")
	}

	tests := []struct {
		name        string
		sessionID   string
		method      string
		expectedMsg string
	}{
		{ // valid
			name:        "Test 1",
			sessionID:   session1.SessionID,
			method:      "POST",
			expectedMsg: "success",
		},
		{ // invalid method
			name:        "Test 2",
			sessionID:   session2.SessionID,
			method:      "GET",
			expectedMsg: api.ErrPostMethod,
		},
		{ // valid
			name:        "Test 3",
			sessionID:   session2.SessionID,
			method:      "POST",
			expectedMsg: "success",
		},
		{ // not logged in anymore
			name:        "Test 4",
			sessionID:   "bruh",
			method:      "POST",
			expectedMsg: api.ErrUnauthorized,
		},
	}

	server := api.NewServer(":3000")
	server.Post("/logout", server.HandleLogout, api.AuthMiddleware)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/logout", nil)
			r.AddCookie(&http.Cookie{
				Name:  api.SESSIONID_COOKIE_NAME,
				Value: tc.sessionID,
			})
			w := httptest.NewRecorder()

			server.Mux.ServeHTTP(w, r)

			msg := strings.TrimSpace(w.Body.String())
			if msg != tc.expectedMsg {
				t.Fatalf("Expected msg: %s, got: %s\n", tc.expectedMsg, msg)
			}

			if tc.method != "POST" {
				return
			}

			_, sessionExists := sessionManager.GetSession(tc.sessionID)
			if sessionExists {
				t.Fatal("Session is not supposed to exist")
			}
		})

	}

}

/*
This test is to see if the server is setting the cookie
for the session in the response header when logging in

This test should be run after TestHandleSignup is called, to ensure
the tested user is added first
*/
func TestLoginCookie(t *testing.T) {
	db.Init()
	defer db.Close()
	server := api.NewServer(":3000")
	server.Post("/login", server.HandleLogin)

	payload := `{"username": "testuser1", "password": "testpassword1"}`
	r := httptest.NewRequest("POST", "/login", bytes.NewBufferString(payload))

	w := httptest.NewRecorder()
	server.Mux.ServeHTTP(w, r)

	cookie := w.Header().Get("Set-Cookie")
	if !strings.Contains(cookie, "afpsid=") {
		t.Fatalf("Cookie has not been set\n")
	}
}

/*
Integration test for HandleListFurniture handler to see if
1. JSON sent by the client is being received by the server correctly
  - From the frontend (React TS), the body is sent as one long string
    representing a JSON object bc of JSON.stringify(obj)

2. New furniture listings are saved in the DB
  - check this by grabbbing the listingID from the response
    and use that to find it in the DB

*From the frontend (React TS), the body is sent as one long string
representing a JSON object

This test uses the sessionID to find the client's document in the "users" collection.
Then it uses the userID from the document to associate the furniture listing with the
client.

Sample sessionID (which belongs to a test acc in the DB) to be used when testing:
sessionID -> "testtest-test-test-test-testtesttest"
*/
func TestHandleListFurniture(t *testing.T) {

	// prepare images to simulate http request from client
	imgData1, err := util.EncodeImageToBase64("../tests/test_images/tiger_maple1.jpg")
	if err != nil {
		t.Fatalf("Error encoding image1 to base64 string")
	}

	imgData2, err := util.EncodeImageToBase64("../tests/test_images/tiger_maple2.jpg")
	if err != nil {
		t.Fatalf("Error encoding image2 to base64 string")
	}

	// creating sample listing
	furnitureListing1 := types.FurnitureListing{
		Title:       "English Tiger maple queen bed",
		Description: "My favorite bed",
		Type:        types.Bed,
		Cost:        7500,
		Style:       types.English,
		Condition:   "Great",
		Material:    types.TigerMaple,
		Images:      []string{imgData1, imgData2},
	}

	furnitureListing2 := types.FurnitureListing{
		Title:       "Cherry Farm Table Sheraton Style",
		Description: "Selling my lovely Cherry Farm Table",
		Type:        types.Table,
		Cost:        2700,
		Style:       types.Sheraton,
		Condition:   "Great",
		Material:    types.Cherry,
	}

	// simulating JSON.stringify(obj) in TS
	payload1, err := json.Marshal(furnitureListing1)
	if err != nil {
		t.Fatal("Failed to encode payload1 into JSON")
	}

	payload2, err := json.Marshal(furnitureListing2)
	if err != nil {
		t.Fatal("Failed to encode payload2 into JSON")
	}

	// creating a fake loggedIn client using the test account
	const TEST_SESS_ID string = "testtest-test-test-test-testtesttest"
	sessionManager := api.GetSessionManager()
	session, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: TEST_SESS_ID,
	})

	if err != nil {
		t.Fatalf("Err creating simulated session: %s\n", err.Error())
	}

	// userID from the test account
	TEST_OBJID, err := primitive.ObjectIDFromHex("65b094f4a2cb3bf5e40d42d7")
	if err != nil {
		t.Fatal(err)
	}

	// when logging in, the userid gets saved into the session store
	session.Store["userid"] = TEST_OBJID

	// create test cases
	tests := []struct {
		name               string
		method             string
		sessionID          string
		payload            string
		expectedStatusCode int
		expectedMessage    string
	}{
		{ // valid
			name:               "Test 1",
			method:             "POST",
			payload:            string(payload1),
			sessionID:          session.SessionID,
			expectedStatusCode: http.StatusOK,
		},
		{ // missing image
			name:               "Test 2",
			method:             "POST",
			payload:            string(payload2),
			sessionID:          session.SessionID,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Images not provided",
		},
		{ // invalid payload format type
			name:               "Test 3",
			method:             "POST",
			payload:            "34",
			sessionID:          session.SessionID,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Failed to decode JSON",
		},
		{ // invalid json formatting
			name:               "Test 4",
			method:             "POST",
			sessionID:          session.SessionID,
			payload:            `{"Title":"Oak Nightstand with refinish}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Failed to decode JSON",
		},
		{ // missing condition
			name:               "Test 5",
			method:             "POST",
			sessionID:          session.SessionID,
			payload:            `{"Title":"Oak Nightstand with refinish"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Condition not provided",
		},
		{ // not logged in
			name:               "Test 6",
			method:             "POST",
			payload:            string(payload1),
			sessionID:          "foo",
			expectedStatusCode: http.StatusUnauthorized,
			expectedMessage:    api.ErrUnauthorized,
		},
		{ // invalid method
			name:               "Test 7",
			method:             "PUT",
			payload:            string(payload1),
			sessionID:          "foo",
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedMessage:    api.ErrPostMethod,
		},
	}

	db.Init()
	defer db.Close()
	server := api.NewServer(":3000")
	server.Post("/list_furniture", server.HandleListFurniture, api.AuthMiddleware)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			r := httptest.NewRequest(tc.method, "/list_furniture", strings.NewReader(tc.payload))
			r.AddCookie(&http.Cookie{
				Name:  api.SESSIONID_COOKIE_NAME,
				Value: tc.sessionID,
			})
			w := httptest.NewRecorder()

			/*
				need to call the http handler after the first error condition bc
				the json decoding of the request body in the handler eats up the
				request body, which would result in an empty body being read if we
				call the http handler before the first check
			*/
			server.Mux.ServeHTTP(w, r)

			status := w.Code
			if status != tc.expectedStatusCode {
				t.Fatalf("Expected status code: %d, got: %v\n", tc.expectedStatusCode, status)
			}

			// validate that body was inserted into MongoDB correctly
			res := strings.TrimSpace(w.Body.String())
			if res == "" {
				t.Fatal("Response did not return an anything")
			}

			// find the document in the listings collection with the listingID and userID
			listingResDB, err := db.FindByIDInListingsCollection(res)

			// compare expected message
			if tc.expectedMessage != "" && res != tc.expectedMessage {
				t.Fatalf("Expected: '%s', got: '%s'\n", tc.expectedMessage, res)
			}
			// stop testcase when an error is reached
			if err != nil {
				return
			}

			/*
				THE CODE BELOW..
				only applies to testcases that are simulating valid requests, like Test 1,
				bc the other test cases, which tests for errors, will return above
			*/

			var actualListing types.FurnitureListing
			decodeErr := listingResDB.Decode(&actualListing)
			if decodeErr != nil {
				t.Fatalf("Failed to decode resultDB into struct: %v\n", decodeErr)
			}

			/*
				use the userID and the listingID returned from the response recorder
				to validate that the listing was added to the DB
			*/

			var objectID primitive.ObjectID = session.Store["userid"].(primitive.ObjectID)
			listingObjectID, err := primitive.ObjectIDFromHex(res)
			if err != nil {
				t.Fatal("Hex string is not a valid objectID")
			}

			listingsColletion := db.GetCollection("listings")
			result := listingsColletion.FindOne(context.Background(), bson.M{
				"userid": objectID,
				"_id":    listingObjectID,
			})
			if result.Err() == mongo.ErrNoDocuments {
				t.Fatalf("Expected document with userID: %s, and listingID: %s; got nothing\n", objectID.Hex(), listingObjectID.Hex())
			}

		})
	}
}

/*
Unit test for ValidateListFormFields
*/
func TestValidateListFormFields(t *testing.T) {
	tests := []struct {
		name     string
		payload  types.FurnitureListing
		expected bool
	}{
		{ // valid FurnitureListing
			name: "Test 1",
			payload: types.FurnitureListing{
				Title:       "something",
				Description: "something",
				Type:        types.Bed,
				Cost:        34.99,
				Style:       "English",
				Condition:   "Great",
				Material:    types.Pine,
				Images:      []string{"1", "2"},
			},
			expected: true,
		},
		{ // missing title
			name: "Test 2",
			payload: types.FurnitureListing{
				Title:       "",
				Description: "something",
				Type:        types.Bed,
				Cost:        34.99,
				Style:       "English",
				Condition:   "Great",
				Material:    types.Pine,
				Images:      []string{"1", "2"},
			},
			expected: false,
		},
		{ // empty values for all fields
			name:     "Test 3",
			payload:  types.FurnitureListing{},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			validated, _ := api.ValidateListFormFields(tc.payload)

			if validated != tc.expected {
				t.Fatalf("Expected: %v, got: %v", tc.expected, validated)
			}
		})
	}
}

func TestHandleGetFurnitures(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		expectedStatusCode int
		expectedMessage    string
	}{
		{ // valid
			name:               "Test 1",
			method:             "GET",
			expectedStatusCode: http.StatusOK,
		},
		{ // valid
			name:               "Test 2",
			method:             "POST",
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedMessage:    "Request must be a GET request",
		},
	}

	db.Init()
	defer db.Close()
	server := api.NewServer(":3000")
	server.Get("/get_furnitures", server.HandleGetFurnitures)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/get_furnitures", nil)
			w := httptest.NewRecorder()

			server.Mux.ServeHTTP(w, r)

			statusCode := w.Code
			if statusCode != tc.expectedStatusCode {
				t.Fatalf("Expected status code: %d, got: %d\n", tc.expectedStatusCode, statusCode)
			}

			if tc.expectedMessage != "" {
				message := strings.TrimSpace(w.Body.String())
				if message != tc.expectedMessage {
					t.Fatalf("Expected msg: %s, got: %s\n", tc.expectedMessage, message)
				}
			}
		})
	}
}

func TestHandleGetFurniture(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		payload            string
		expectedStatusCode int
		expectedMessage    string
	}{
		{ // invalid ID
			name:               "Test 1",
			method:             "GET",
			payload:            "",
			expectedMessage:    "Furniture listing with provided listingID not found",
			expectedStatusCode: http.StatusBadRequest,
		},
		{ // valid
			name:               "Test 2",
			method:             "GET",
			expectedMessage:    "success",
			payload:            "65b433dd8d3c8f926b88cd7a",
			expectedStatusCode: http.StatusOK,
		},
		{ // invalid method
			name:               "Test 3",
			method:             "POST",
			expectedMessage:    "Request must be a GET request",
			payload:            "234324324234324",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{ // invalid method
			name:               "Test 4",
			method:             "DELETE",
			expectedMessage:    "Request must be a GET request",
			payload:            "6783",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{ // invalid listingID
			name:               "Test 5",
			method:             "GET",
			expectedMessage:    "Furniture listing with provided listingID not found",
			payload:            "65aeadc6f96ba92452b8e52",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	db.Init()
	defer db.Close()
	server := api.NewServer(":3000")
	server.Get("/get_furniture", server.HandleGetFurniture)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var target string = fmt.Sprintf("/get_furniture?listingid=%s", tc.payload)
			r := httptest.NewRequest(tc.method, target, nil)
			w := httptest.NewRecorder()

			server.Mux.ServeHTTP(w, r)

			if strings.TrimSpace(w.Body.String()) != tc.expectedMessage {
				t.Fatalf("Expected msg: %s, got: %s\n", tc.expectedMessage, w.Body.String())
			}

			if w.Code != tc.expectedStatusCode {
				t.Fatalf("Expected status code: %d, got: %d\n", tc.expectedStatusCode, w.Code)
			}

		})
	}
}

// 1/28 [WIP]
// func TestHandleAccountGet(t *testing.T) {
// 	// creating fake logged in clients
// 	sessionManager := api.GetSessionManager()
// 	session1, err := sessionManager.CreateSession(api.SessionTemplate{
// 		SessionID: "",
// 	})
// 	if err != nil {
// 		t.Fatal("Failed to generate session1")
// 	}
// 	session1.Store["userid"] = "testuser1"

// 	session2, err := sessionManager.CreateSession(api.SessionTemplate{
// 		SessionID: "",
// 	})
// 	if err != nil {
// 		t.Fatal("Failed to generate session1")
// 	}
// 	session2.Store["userid"] = "testuser2"

// 	tests := []struct {
// 		name        string
// 		method      string
// 		sessionid   string
// 		expectedMsg string
// 	}{
// 		{ // valid, authorized
// 			name:        "Test 1",
// 			method:      "GET",
// 			sessionid:   session1.SessionID,
// 			expectedMsg: "success",
// 		},
// 		{ // valid, authorized
// 			name:        "Test 2",
// 			method:      "GET",
// 			sessionid:   session2.SessionID,
// 			expectedMsg: "success",
// 		},
// 		{ // unauthorized
// 			name:        "Test 3",
// 			method:      "GET",
// 			sessionid:   "foo",
// 			expectedMsg: api.ErrUnauthorized,
// 		},
// 	}

// 	server := api.NewServer(":3000")
// 	server.Use("/account", server.HandleAccount, api.AuthMiddleware)

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			r := httptest.NewRequest(tc.method, "/account", nil)
// 			r.AddCookie(&http.Cookie{
// 				Name:  api.SESSIONID_COOKIE_NAME,
// 				Value: tc.sessionid,
// 			})
// 			w := httptest.NewRecorder()

// 			server.Mux.ServeHTTP(w, r)

// 			res := strings.TrimSpace(w.Body.String())

// 			if res != tc.expectedMsg {
// 				t.Fatalf("Expected msg: %s, got: %s\n", tc.expectedMsg, res)
// 			}
// 		})
// 	}
// }

/*
This test is WIP
*/
func TestHandleCheckout(t *testing.T) {
	// creating fake loggedIn user with test account
	sessionManager := api.GetSessionManager()
	session1, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "testtest-test-test-test-testtesttest",
	})
	if err != nil {
		t.Fatal("Failed to create a fake session")
	}

	USER_ID, err := primitive.ObjectIDFromHex("65b094f4a2cb3bf5e40d42d7")
	if err != nil {
		t.Fatal("Failed to generate objectID for fake session")
	}
	session1.Store["userid"] = USER_ID

	// mock checkout data
	checkoutInfo := api.CheckoutInfo{
		ShoppingCart: []api.FurnitureItem{
			"65b433dd8d3c8f926b88cd7a",
		},
		Payment: api.PaymentInfo{
			StripeToken:   "token_foo",
			PaymentMethod: "Credit",
			Amount:        7500,
			Currency:      "usd",
		},
		ShippingAddress: api.ShippingAddress{
			State:   "RI",
			City:    "Providence",
			Street:  "999 Holy St.",
			ZipCode: "02907",
		},
	}

	// mock expected receipt
	expectedReceipt := api.Receipt{
		ShippingAddress: checkoutInfo.ShippingAddress,
		PaymentMethod:   "Credit",
		TotalCost:       7500,
		Items:           []string{string(checkoutInfo.ShoppingCart[0])},
		UserID:          session1.Store["userid"].(primitive.ObjectID),
	}

	// encode checkoutInfo
	checkoutJSONData, err := json.Marshal(checkoutInfo)
	if err != nil {
		t.Fatal("Failed to encode checkoutInfo into JSON")
	}

	// encode expectedReceipt
	receiptJSONData, err := json.Marshal(expectedReceipt)
	if err != nil {
		t.Fatal("Failed to encode expectedReceipt into JSON")
	}

	tests := []struct {
		name               string
		method             string
		sessionid          string
		payload            string
		expectedStatusCode int
		expectedMsg        string
	}{
		{ // unauthorized user
			name:               "Test 1",
			method:             "POST",
			sessionid:          "unauthorized",
			payload:            string(checkoutJSONData),
			expectedStatusCode: http.StatusUnauthorized,
			expectedMsg:        api.ErrUnauthorized,
		},
		{ // invalid method
			name:               "Test 2",
			method:             "GET",
			sessionid:          session1.SessionID,
			payload:            string(checkoutJSONData),
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedMsg:        api.ErrPostMethod,
		},
		{ // valid
			name:               "Test 3",
			method:             "POST",
			sessionid:          session1.SessionID,
			payload:            string(checkoutJSONData),
			expectedStatusCode: http.StatusOK,
			expectedMsg:        string(receiptJSONData),
		},
	}

	db.Init()
	defer db.Close()
	server := api.NewServer(":3000")
	server.Post("/checkout", server.HandleCheckout, api.AuthMiddleware)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/checkout", strings.NewReader(tc.payload))
			r.AddCookie(&http.Cookie{
				Name:  api.SESSIONID_COOKIE_NAME,
				Value: tc.sessionid,
			})
			w := httptest.NewRecorder()

			server.Mux.ServeHTTP(w, r)

			res := strings.TrimSpace(w.Body.String())

			if w.Code != tc.expectedStatusCode {
				t.Fatalf("Expected code: %d, got: %d\n", tc.expectedStatusCode, w.Code)
			}

			if res != tc.expectedMsg {
				t.Fatalf("Expected msg: %s, got: %s\n", tc.expectedMsg, res)
			}
		})
	}
}
