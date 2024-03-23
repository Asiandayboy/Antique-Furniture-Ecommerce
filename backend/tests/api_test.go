package tests

import (
	"backend/api"
	"backend/db"
	"backend/types"
	"backend/util"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"

	// "backend/util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
	server.Use("POST /signup", server.HandleSignup)

	tests := []struct {
		name               string
		method             string
		payload            string
		expectedResMsg     string
		expectedStatusCode int
	}{
		// { // test valid signup; if this test fails, it's bc the payload is already in the db; use different data
		// 	name:               "Test 1",
		// 	method:             "POST",
		// 	payload:            `{"username": "testuser1", "password": "testpassword1", "confirm": "testpassword1", "email": "test@gmail.com"}`,
		// 	expectedResMsg:     "success",
		// 	expectedStatusCode: http.StatusOK,
		// },
		{ // test invalid method
			name:               "Test 2",
			method:             "GET",
			payload:            `{"username": "testuser1", "password": "testpassword1", "confirm": "testpassword1", "email": "test@gmail.com"}`,
			expectedResMsg:     api.ErrMethodNotAllowed,
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{ // testing invalid json decode
			name:               "Test 3",
			method:             "POST",
			payload:            `{"username": "testuser1", "password": "testpassword1, "confirm": "testpassword1", "email": "test@gmail.com"}`,
			expectedResMsg:     "Could not decode request body into JSON",
			expectedStatusCode: http.StatusBadRequest,
		},
		{ // testing username against existing username -> "bob"
			name:               "Test 4",
			method:             "POST",
			payload:            `{"username": "bob", "password": "testpassword1", "confirm": "testpassword1", "email": "test@gmail.com"}`,
			expectedResMsg:     "Username is taken",
			expectedStatusCode: http.StatusConflict,
		},
		{ // testing email against existing email -> "johnsmith@gmail.com"
			name:               "Test 5",
			method:             "POST",
			payload:            `{"username": "testuser123", "password": "testpassword1", "confirm": "testpassword1", "email": "johnsmith@gmail.com"}`,
			expectedResMsg:     "Email is taken",
			expectedStatusCode: http.StatusConflict,
		},
		{ // testing password mismatch
			name:               "Test 6",
			method:             "POST",
			payload:            `{"username": "testuser123", "password": "password", "confirm": "testpassword1", "email": "johnsmith@gmail.com"}`,
			expectedResMsg:     api.ErrPasswordMismatch,
			expectedStatusCode: http.StatusBadRequest,
		},
		{ // testing empty fields
			name:               "Test 7",
			method:             "POST",
			payload:            `{"username": "testuser123", "password": "password"}`,
			expectedResMsg:     api.ErrBlankFields,
			expectedStatusCode: http.StatusBadRequest,
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
	server.Use("POST /login", server.HandleLogin)

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
			expectedResMsg:     api.ErrMethodNotAllowed,
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
			expectedMsg: api.ErrMethodNotAllowed,
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
	server.Use("POST /logout", server.HandleLogout, api.AuthMiddleware)
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
	server.Use("POST /login", server.HandleLogin)

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

	/*------------------Test Data 1 (valid)-------------------*/

	furnitureListing1 := types.FurnitureListing{
		Title:       "Oak nightstand double drawers",
		Description: "Fit all your personal belongings",
		Type:        types.Nightstand,
		Cost:        1000,
		Style:       types.Federal,
		Condition:   types.Good,
		Material:    types.Oak,
	}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// add JSON data to multipart writer
	jsonPart, err := writer.CreateFormField("json_data")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewEncoder(jsonPart).Encode(furnitureListing1); err != nil {
		t.Fatal(err)
	}

	// add images to multipart writer
	filePaths := []string{"../tests/test_images/tiger_maple1.jpg", "../tests/test_images/tiger_maple2.jpg"}

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			t.Fatal("Failed to open file", err)
		}
		defer file.Close()

		filePart, err := writer.CreateFormFile("furniture_images", filepath.Base(filePath))
		if err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(filePart, file); err != nil {
			t.Fatal(err)
		}
	}

	writer.Close()

	/*------------------Test Data 2 (no images provided)-------------------*/

	furnitureListing2 := types.FurnitureListing{
		Title:       "Cherry Farm Table Sheraton Style",
		Description: "Selling my lovely Cherry Farm Table",
		Type:        types.Table,
		Cost:        2700,
		Style:       types.Sheraton,
		Condition:   "Great",
		Material:    types.Cherry,
	}

	var requestBody2 bytes.Buffer
	writer2 := multipart.NewWriter(&requestBody2)

	jsonPart2, err := writer2.CreateFormField("json_data")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewEncoder(jsonPart2).Encode(furnitureListing2); err != nil {
		t.Fatal(err)
	}

	writer2.Close()

	/*------------------Test Data 3 (all list form fields empty)-------------------*/

	var requestBody3 bytes.Buffer
	writer3 := multipart.NewWriter(&requestBody3)

	jsonPart3, err := writer3.CreateFormField("json_data")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewEncoder(jsonPart3).Encode(types.FurnitureListing{}); err != nil {
		t.Fatal(err)
	}

	writer3.Close()

	/*-----------------Create fake logged in user-----------------*/

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

	/*-------------------------------------------------------------*/

	// create test cases
	tests := []struct {
		name               string
		method             string
		sessionID          string
		payload            *bytes.Buffer
		writer             *multipart.Writer
		expectedStatusCode int
		expectedMessage    string
	}{
		{ // valid
			name:               "Test 1",
			method:             "POST",
			payload:            &requestBody,
			writer:             writer,
			sessionID:          session.SessionID,
			expectedStatusCode: http.StatusOK,
		},
		{ // missing image
			name:               "Test 2",
			method:             "POST",
			payload:            &requestBody2,
			writer:             writer2,
			sessionID:          session.SessionID,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    `["Furniture images not provided"]`,
		},
		{ // every field missing
			name:               "Test 3",
			method:             "POST",
			sessionID:          session.SessionID,
			payload:            &requestBody3,
			writer:             writer3,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    api.ErrListFormEveryFieldMissing,
		},
		{ // not logged in
			name:               "Test 4",
			method:             "POST",
			payload:            &requestBody,
			writer:             writer,
			sessionID:          "foo",
			expectedStatusCode: http.StatusUnauthorized,
			expectedMessage:    api.ErrUnauthorized,
		},
		{ // invalid method
			name:               "Test 5",
			method:             "PUT",
			payload:            &requestBody,
			writer:             writer,
			sessionID:          "foo",
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedMessage:    api.ErrMethodNotAllowed,
		},
	}

	db.Init()
	defer db.Close()
	server := api.NewServer(":3000")
	server.Use("POST /list_furniture", server.HandleListFurniture, api.AuthMiddleware)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/list_furniture", tc.payload)
			r.Header.Add("Content-Type", tc.writer.FormDataContentType())
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

			// validate that body was inserted into MongoDB correctly
			res := strings.TrimSpace(w.Body.String())
			if res == "" {
				t.Fatal("Response did not return an anything")
			}

			status := w.Code
			if status != tc.expectedStatusCode {
				t.Fatalf("Expected status code: %d, got: %v\n", tc.expectedStatusCode, status)
			}

			// find the document in the listings collection with the listingID and userID
			listingResDB, err := db.FindByIDInListingsCollection(res)

			// compare expected message
			if tc.expectedMessage != "" && res != tc.expectedMessage {
				t.Fatalf("Expected msg: '%s', got: '%s'\n", tc.expectedMessage, res)
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

			var objectID primitive.ObjectID = actualListing.UserID
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
	payload1 := types.FurnitureListing{
		Title:       "something",
		Description: "something",
		Type:        types.Bed,
		Cost:        34.99,
		Style:       "English",
		Condition:   "Great",
		Material:    types.Pine,
		// Images:      []string{"1", "2"},
	}

	payload2 := types.FurnitureListing{
		Title:       "",
		Description: "something",
		Type:        types.Bed,
		Cost:        34.99,
		Style:       "English",
		Material:    types.Pine,
		// Images:      []string{"1", "2"},
	}

	payload3 := types.FurnitureListing{}

	tests := []struct {
		name           string
		payload        types.FurnitureListing
		expected       bool
		expectedErrMsg string
	}{
		{ // valid FurnitureListing
			name:     "Test 1",
			payload:  payload1,
			expected: true,
		},
		{ // missing title
			name:           "Test 2",
			payload:        payload2,
			expected:       false,
			expectedErrMsg: `["Furniture condition not provided","Furniture title not provided"]`,
		},
		{ // empty values for all fields
			name:           "Test 3",
			payload:        payload3,
			expected:       false,
			expectedErrMsg: api.ErrListFormEveryFieldMissing,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := api.ValidateListFormFields(tc.payload)

			if tc.expected {
				if err != nil {
					t.Fatalf("Expected: %v, got: %v", tc.expected, err)
				}
			} else {
				if err.Error() != tc.expectedErrMsg {
					t.Fatalf("Expected: %s, got: %s", tc.expectedErrMsg, err.Error())
				}
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
			expectedMessage:    api.ErrMethodNotAllowed,
		},
	}

	db.Init()
	defer db.Close()
	server := api.NewServer(":3000")
	server.Use("GET /get_furnitures", server.HandleGetFurnitures)
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
	}{
		{ // invalid ID
			name:               "Test 1",
			method:             "GET",
			payload:            "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{ // valid
			name:               "Test 2",
			method:             "GET",
			payload:            "65bf607585af14e593096ea1",
			expectedStatusCode: http.StatusOK,
		},
		{ // invalid method
			name:               "Test 3",
			method:             "POST",
			payload:            "234324324234324",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{ // invalid method
			name:               "Test 4",
			method:             "DELETE",
			payload:            "6783",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{ // invalid listingID
			name:               "Test 5",
			method:             "GET",
			payload:            "65aeadc6f96ba92452b8e52",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	db.Init()
	defer db.Close()
	server := api.NewServer(":3000")
	server.Use("GET /get_furniture", server.HandleGetFurniture)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var target string = fmt.Sprintf("/get_furniture?listingid=%s", tc.payload)
			r := httptest.NewRequest(tc.method, target, nil)
			w := httptest.NewRecorder()

			server.Mux.ServeHTTP(w, r)

			if w.Code != tc.expectedStatusCode {
				t.Fatalf("Expected status code: %d, got: %d\n", tc.expectedStatusCode, w.Code)
			}

		})
	}
}

func TestHandleAccountGET(t *testing.T) {
	db.Init()
	defer db.Close()
	usersCollection := db.GetCollection("users")
	sessionManager := api.GetSessionManager()

	/*-----------------Fake logged in user 1-----------------*/

	session1, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "testtest-test-test-test-testtesttest",
	})
	if err != nil {
		t.Fatal("Failed to generate session1")
	}

	objID1, _ := primitive.ObjectIDFromHex("65b094f4a2cb3bf5e40d42d7")
	session1.Store["userid"] = objID1

	// prepare expected data
	var testuser1ExpectedData types.User
	err = usersCollection.FindOne(context.Background(), bson.M{
		"_id": objID1,
	}).Decode(&testuser1ExpectedData)

	if err != nil {
		t.Fatal(err)
	}

	expectedJSON1, err := json.Marshal(testuser1ExpectedData)
	if err != nil {
		t.Fatal("Failed to marshal testdata1")
	}

	/*----------------Fake logged in user 2------------------*/

	session2, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "bobbobbo-bobb-bobb-bobb-bobbobbobbob",
	})
	if err != nil {
		t.Fatal("Failed to generate session1")
	}

	objID2, _ := primitive.ObjectIDFromHex("65a5a07f062510f606cbd0ae")
	session2.Store["userid"] = objID2

	// prepare expected data
	var testuser2ExpectedData types.User
	err = usersCollection.FindOne(context.Background(), bson.M{
		"_id": objID2,
	}).Decode(&testuser2ExpectedData)

	if err != nil {
		t.Fatal(err)
	}

	expectedJSON2, err := json.Marshal(testuser2ExpectedData)
	if err != nil {
		t.Fatal("Failed to marshal testdata2")
	}

	/*-----------------------test cases------------------------*/

	tests := []struct {
		name               string
		method             string
		sessionid          string
		expectedStatusCode int
		expectedMsg        string
	}{
		{ // valid, authorized
			name:               "Test 1",
			method:             "GET",
			sessionid:          session1.SessionID,
			expectedStatusCode: http.StatusOK,
			expectedMsg:        string(expectedJSON1),
		},
		{ // valid, authorized
			name:               "Test 2",
			method:             "GET",
			sessionid:          session2.SessionID,
			expectedStatusCode: http.StatusOK,
			expectedMsg:        string(expectedJSON2),
		},
		{ // unauthorized
			name:               "Test 3",
			method:             "GET",
			sessionid:          "foo",
			expectedStatusCode: http.StatusUnauthorized,
			expectedMsg:        api.ErrUnauthorized,
		},
	}

	server := api.NewServer(":3000")
	server.Use("/account", server.HandleAccountGET, api.AuthMiddleware)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/account", nil)
			r.AddCookie(&http.Cookie{
				Name:  api.SESSIONID_COOKIE_NAME,
				Value: tc.sessionid,
			})
			w := httptest.NewRecorder()

			server.Mux.ServeHTTP(w, r)

			res := strings.TrimSpace(w.Body.String())

			if w.Code != tc.expectedStatusCode {
				t.Fatalf("Expected status code: %d, got: %d\n", tc.expectedStatusCode, w.Code)
			}

			if res != tc.expectedMsg {
				t.Fatalf("Expected msg: %s, got: %s\n", tc.expectedMsg, res)
			}
		})
	}
}

/*
Be sure to change the payload for each testcase to something different
each time you run the test. Or else, you really won't be testing anything
because the data will stay the same, which is to be expected. But, that's
not what we're testing for, ofc.
*/
func TestHandleAccountPUT(t *testing.T) {
	db.Init()
	defer db.Close()
	// usersCollection := db.GetCollection("users")
	sessionManager := api.GetSessionManager()

	/*---------------fake logged in user1------------------*/

	session1, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "testtest-test-test-test-testtesttest",
	})
	if err != nil {
		t.Fatal("Failed to generate session1")
	}

	objID1, _ := primitive.ObjectIDFromHex("65b094f4a2cb3bf5e40d42d7")
	session1.Store["userid"] = objID1

	// change payload before you run the test
	accountEdit1 := api.AccountEdit{
		NewPhone:    "101-111-4444",
		NewPassword: "test1234",
	}

	payload1, err := json.Marshal(accountEdit1)
	if err != nil {
		t.Fatal("Failed to marshal payload1")
	}

	/*---------------------testcases----------------------*/

	tests := []struct {
		name               string
		sessionID          string
		userID             primitive.ObjectID
		payload            string
		expectedStatusCode int
		expectedMsg        string
	}{
		{
			name:               "Test 1",
			sessionID:          session1.SessionID,
			userID:             objID1,
			payload:            string(payload1),
			expectedStatusCode: http.StatusOK,
			expectedMsg:        "success",
		},
	}

	server := api.NewServer(":3000")
	server.Use("PUT /account", server.HandleAccountPUT, api.AuthMiddleware)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("PUT", "/account", strings.NewReader(tc.payload))
			r.AddCookie(&http.Cookie{
				Name:  api.SESSIONID_COOKIE_NAME,
				Value: tc.sessionID,
			})
			w := httptest.NewRecorder()

			server.Mux.ServeHTTP(w, r)

			res := strings.TrimSpace(w.Body.String())

			if w.Code != tc.expectedStatusCode {
				t.Fatalf("Expected status code: %d, got: %d\n", tc.expectedStatusCode, w.Code)
			}

			if res != tc.expectedMsg {
				t.Fatalf("Expected msg: %s, got: %s\n", tc.expectedMsg, res)
			}

			/*----------------check db to see if the changes went through--------------*/

			usersCollection := db.GetCollection("users")
			var user types.User
			err := usersCollection.FindOne(context.Background(), bson.M{
				"_id": tc.userID,
			}).Decode(&user)

			if err != nil {
				t.Fatal("Failed to find user document for test")
			}

			var changes api.AccountEdit
			err = json.Unmarshal([]byte(tc.payload), &changes)
			if err != nil {
				t.Fatal("Failed to unmarshal payload for test")
			}

			if changes.NewEmail != "" && changes.NewEmail != user.Email {
				t.Fatalf("Expected email: %s, got: %s\n", changes.NewEmail, user.Email)
			}

			if changes.NewPhone != "" && changes.NewPhone != user.Phone {
				t.Fatalf("Expected phone: %s, got: %s\n", changes.NewPhone, user.Phone)
			}

			if changes.NewPassword != "" {
				// compare the plain text of the new password and the resulting hash of it
				err = util.CheckPassword(changes.NewPassword, user.Password)
				if err != nil {
					t.Fatal("Passwords are not the same")
				}
			}
		})
	}

}

/*
This test validates that the client receives the expected
receipt. This test is currently simulating a fake user
initiating a checkout session

This test is merely being used to see if a checkout session is
being created, if it triggers the webhook endpoint, and if it can
retrieve the session data that was passed as metadata; it's not
really a conventional test.

THIS TEST IS NOT REALLY A TEST; I JUST WANNA SEE IF I INTEGRATED
STRIPE CORRECTLY. The other way I'll test this is when I create my
frontend

# Test card used to send test funds from test transactions to platform account balance

Find here and scroll down to Available balance ---> https://stripe.com/docs/testing

- number: 4000000000000077
- CVC: any 3 digit number
- Expiration Date : any date in the future
- ZipCode: any 5 digit number

All funds go to the platform account. Then when users go to their dashboard,
look at their funds that they've received from selling their antique furnitures,
they can press a "payout" button, enter the amount they have available, and it
will take from the platform account the funds and transfer it to their bank
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

	// mock checkout input data
	checkoutInfo := api.CheckoutInfo{
		ShoppingCart: []string{
			"65bf607585af14e593096ea1",
		},
		Payment: api.PaymentInfo{
			PaymentMethod: "Credit",
			Amount:        7500,
			Currency:      "usd",
		},
	}

	checkoutJSONData, err := json.Marshal(checkoutInfo)
	if err != nil {
		t.Fatal("Failed to encode checkoutInfo into JSON")
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
			expectedMsg:        api.ErrMethodNotAllowed,
		},
		{ // valid [*This test doesn't pass yet]
			name:               "Test 3",
			method:             "POST",
			sessionid:          session1.SessionID,
			payload:            string(checkoutJSONData),
			expectedStatusCode: http.StatusOK,
			expectedMsg:        "",
		},
	}

	db.Init()
	defer db.Close()
	server := api.NewServer(":3000")
	go server.Start()

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

			if res != tc.expectedMsg && tc.expectedMsg != "" {
				t.Fatalf("Expected msg: %s, got: %s\n", tc.expectedMsg, res)
			}
		})
	}

	for {
		time.Sleep(time.Second)
	}
}

func TestHandleAddressGET(t *testing.T) {
	db.Init()
	defer db.Close()
	sessionManager := api.GetSessionManager()

	/*-----------Fake logged in user 1-------------*/

	session1, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "testtest-test-test-test-testtesttest",
	})
	if err != nil {
		t.Fatal("Failed to create session for testuser1")
	}
	userID1, _ := primitive.ObjectIDFromHex("65b094f4a2cb3bf5e40d42d7")
	session1.Store["userid"] = userID1

	// generating expected msg
	var expectedAddresses1 []types.ShippingAddress
	usersCollection := db.GetCollection("shippingAddresses")
	cursor, err := usersCollection.Find(
		context.Background(),
		bson.M{"userid": userID1},
	)
	if err != nil {
		t.Fatal("Failed to find documents")
	}
	err = cursor.All(context.Background(), &expectedAddresses1)
	if err != nil {
		t.Fatal(err.Error())
	}

	jsonData1, err := json.Marshal(expectedAddresses1)
	if err != nil {
		t.Fatal("Failed to marshal expectedAddresses1")
	}

	expectedMsg1 := string(jsonData1)

	/*-----------------test cases------------------*/

	tests := []struct {
		name               string
		sessionID          string
		expectedStatusCode int
		expectedMsg        string
	}{
		{ // valid
			name:               "Test 1",
			sessionID:          session1.SessionID,
			expectedStatusCode: http.StatusOK,
			expectedMsg:        expectedMsg1,
		},
		{ // unauthorized
			name:               "Test 2",
			sessionID:          "notloggedin",
			expectedStatusCode: http.StatusUnauthorized,
			expectedMsg:        api.ErrUnauthorized,
		},
	}

	server := api.NewServer(":3000")
	server.Use("GET /account/address", server.HandleAddressGET, api.AuthMiddleware)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/account/address", nil)
			r.AddCookie(&http.Cookie{
				Name:  api.SESSIONID_COOKIE_NAME,
				Value: tc.sessionID,
			})
			w := httptest.NewRecorder()

			server.Mux.ServeHTTP(w, r)

			res := strings.TrimSpace(w.Body.String())
			code := w.Code

			if code != tc.expectedStatusCode {
				t.Fatalf("Expected status code: %d, got: %d\n", tc.expectedStatusCode, code)
			}

			if res != tc.expectedMsg {
				t.Fatalf("Expected msg: %s, got: %s\n", tc.expectedMsg, res)
			}
		})
	}
}

func TestHandleAddressPOST(t *testing.T) {
	db.Init()
	defer db.Close()
	sessionManager := api.GetSessionManager()

	/*-----------Fake logged in user 1-------------*/

	session1, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "testtest-test-test-test-testtesttest",
	})
	if err != nil {
		t.Fatal("Failed to create session for testuser1")
	}
	userID1, _ := primitive.ObjectIDFromHex("65b094f4a2cb3bf5e40d42d7")
	session1.Store["userid"] = userID1

	// addressID1, _ := primitive.ObjectIDFromHex("abcdefabcdefabcdefabcde1")
	var testInputAddress1 types.ShippingAddress = types.ShippingAddress{
		// AddressID: addressID1,
		// UserID:    userID1,
		State:   "RI",
		City:    "Providence",
		Street:  "23 Something Street",
		ZipCode: "11111",
		Default: false,
	}

	jsonData1, err := json.Marshal(testInputAddress1)
	if err != nil {
		t.Fatal("Failed to marshal address1")
	}

	/*-----------------test cases------------------*/

	/*
		Should expect the address to be returned in JSON and
		match the one that was inserted into the DB
	*/

	tests := []struct {
		name               string
		sessionID          string
		payload            string
		expectedStatusCode int
		expectedMsg        string
	}{
		{
			name:               "Test 1",
			sessionID:          session1.SessionID,
			payload:            string(jsonData1),
			expectedStatusCode: http.StatusOK,
			expectedMsg:        "success",
		},
	}

	server := api.NewServer(":3000")
	server.Use("POST /account/address", server.HandleAddressPOST, api.AuthMiddleware)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/account/address", strings.NewReader(tc.payload))
			r.AddCookie(&http.Cookie{
				Name:  api.SESSIONID_COOKIE_NAME,
				Value: tc.sessionID,
			})
			w := httptest.NewRecorder()

			server.Mux.ServeHTTP(w, r)

			res := strings.TrimSpace(w.Body.String())
			code := w.Code

			if code != tc.expectedStatusCode {
				t.Fatalf("Expected status code: %d, got: %d\n", tc.expectedStatusCode, code)
			}

			if res != tc.expectedMsg {
				t.Fatalf("Expected msg: %s, got: %s\n", tc.expectedMsg, res)
			}
		})
	}
}

/*
Be sure to change all the inputs to something different
every time you run this test if you want to see the changes

The AddressID should stay the same to make this test
simpler by just checking for changes in the same document
*/
func TestHandleAddressPUT(t *testing.T) {
	db.Init()
	defer db.Close()
	sessionManager := api.GetSessionManager()

	/*-----------Fake logged in user 1-------------*/

	session1, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "testtest-test-test-test-testtesttest",
	})
	if err != nil {
		t.Fatal("Failed to create session for testuser1")
	}
	userID1, _ := primitive.ObjectIDFromHex("65b094f4a2cb3bf5e40d42d7")
	session1.Store["userid"] = userID1

	/*-----------------test cases------------------*/

	// valid input
	input1 := api.AddressUpdateInput{
		AddressID: "65c3f9c52dd8587a26714756",
		Changes: api.AddressUpdate{
			NewStreet: "105 Wizard Avenue",
		},
	}
	jsonData1, err := json.Marshal(input1)
	if err != nil {
		t.Fatal("Failed to marshal test input1")
	}

	// valid input
	input2 := api.AddressUpdateInput{
		AddressID: "65c3f9c52dd8587a26714756",
		Changes: api.AddressUpdate{
			NewZipCode: "02907",
			NewCity:    "Providence",
			NewState:   "RI",
		},
	}
	jsonData2, err := json.Marshal(input2)
	if err != nil {
		t.Fatal("Failed to marshal test input1")
	}

	// empty input
	input3 := api.AddressUpdateInput{
		AddressID: "65c3f9c52dd8587a26714756",
		Changes:   api.AddressUpdate{},
	}
	jsonData3, err := json.Marshal(input3)
	if err != nil {
		t.Fatal("Failed to marshal test input1")
	}

	tests := []struct {
		name               string
		sessionID          string
		payload            string
		expectedStatusCode int
		expectedMsg        string
	}{
		{
			name:               "Test 1",
			sessionID:          session1.SessionID,
			payload:            string(jsonData1),
			expectedStatusCode: http.StatusOK,
			expectedMsg:        "success",
		},
		{
			name:               "Test 2",
			sessionID:          session1.SessionID,
			payload:            string(jsonData2),
			expectedStatusCode: http.StatusOK,
			expectedMsg:        "success",
		},
		{
			name:               "Test 3",
			sessionID:          session1.SessionID,
			payload:            string(jsonData3),
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        api.ErrAddressNoChanges,
		},
	}

	server := api.NewServer(":3000")
	server.Use("PUT /account/address", server.HandleAddressPUT, api.AuthMiddleware)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("PUT", "/account/address", strings.NewReader(tc.payload))
			r.AddCookie(&http.Cookie{
				Name:  api.SESSIONID_COOKIE_NAME,
				Value: tc.sessionID,
			})
			w := httptest.NewRecorder()

			server.Mux.ServeHTTP(w, r)

			code := w.Code
			resMsg := strings.TrimSpace(w.Body.String())

			if resMsg != tc.expectedMsg {
				t.Fatalf("Expected msg: %s, got: %s\n", tc.expectedMsg, resMsg)
			}
			if code != tc.expectedStatusCode {
				t.Fatalf("Expected status code: %d, got: %d\n", tc.expectedStatusCode, code)
			}

			/*------------compare with DB--------------*/

			/*
				Compare the expected changes with the current state of the
				document in the database to ensure the changes actually saved
			*/

			var inputChanges api.AddressUpdateInput
			err := json.Unmarshal([]byte(tc.payload), &inputChanges)
			if err != nil {
				t.Fatal("Failed to unmarshal payload into <inputChanges>")
			}

			shippingAddrCollection := db.GetCollection("shippingAddresses")

			inputKey := reflect.TypeOf(inputChanges.Changes)
			inputValue := reflect.ValueOf(inputChanges.Changes)
			for i := 0; i < inputKey.NumField(); i++ {
				field := inputKey.Field(i)
				value := inputValue.Field(i)

				// find the changed value in the DB to confirm changes
				if !value.IsZero() {
					// remove "New" and transform to lower to match field name in DB
					cleanedStr1 := strings.Replace(field.Name, "New", "", 1)

					/*
						only transform the first uppercase char after removing "new"
						to lower case to maintain camel case

						NewZipCode -> ZipCode -> zipCode
						NewStreet -> Street -> street
					*/
					field.Name = strings.ToLower(string(cleanedStr1[0])) + cleanedStr1[1:]

					id, _ := primitive.ObjectIDFromHex(inputChanges.AddressID)
					res := shippingAddrCollection.FindOne(
						context.Background(),
						bson.M{"_id": id, field.Name: value.String()},
					)
					if res.Err() != nil {
						t.Fatal("Failed to find document; failed to confirm changes")
					}
				}
			}

		})
	}
}

func TestHandleAddressDELETE(t *testing.T) {
	db.Init()
	defer db.Close()
	sessionManager := api.GetSessionManager()
	/*-----------Fake logged in user 1-------------*/

	session1, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "testtest-test-test-test-testtesttest",
	})
	if err != nil {
		t.Fatal("Failed to create session for testuser1")
	}
	userID1, _ := primitive.ObjectIDFromHex("65b094f4a2cb3bf5e40d42d7")
	session1.Store["userid"] = userID1

	/*----------------------tests--------------------*/

	// create a dummy address document to delete
	shippingAddrColl := db.GetCollection("shippingAddresses")
	mockAddr := types.ShippingAddress{
		UserID:  userID1,
		State:   "RI",
		City:    "Providence",
		Street:  "111 Urmom Avenue",
		ZipCode: "02905",
		Default: false,
	}
	res, err := shippingAddrColl.InsertOne(context.Background(), mockAddr)
	if err != nil {
		t.Fatal("Failed to insert mock document")
	}

	tests := []struct {
		name               string
		sessionID          string
		method             string
		shippAddrID        string
		expectedStatusCode int
		expectedMsg        string
	}{
		{ // valid
			name:               "Test 1",
			sessionID:          session1.SessionID,
			method:             "DELETE",
			shippAddrID:        res.InsertedID.(primitive.ObjectID).Hex(),
			expectedStatusCode: http.StatusOK,
			expectedMsg:        "success",
		},
		{ // invalid ID
			name:               "Test 2",
			sessionID:          session1.SessionID,
			method:             "DELETE",
			shippAddrID:        "testing-testing",
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        primitive.ErrInvalidHex.Error(),
		},
		{ // method not allowed
			name:               "Test 3",
			sessionID:          session1.SessionID,
			method:             "POST",
			shippAddrID:        "testing-testing",
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedMsg:        api.ErrMethodNotAllowed,
		},
	}

	server := api.NewServer(":3000")
	server.Use("DELETE /account/address/{addressID}", server.HandleAddressDELETE, api.AuthMiddleware)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			target := fmt.Sprintf("/account/address/%s", tc.shippAddrID)
			r := httptest.NewRequest(tc.method, target, nil)
			r.AddCookie(&http.Cookie{
				Name:  api.SESSIONID_COOKIE_NAME,
				Value: tc.sessionID,
			})
			w := httptest.NewRecorder()

			/*
				Ensure mock document was inserted for valid request tests
			*/
			id, _ := primitive.ObjectIDFromHex(tc.shippAddrID)
			if tc.expectedMsg == "success" {
				res := shippingAddrColl.FindOne(
					context.Background(),
					bson.M{"_id": id},
				)
				if res.Err() != nil {
					t.Fatal("Document did not get insert into DB properly")
				}
			}

			server.Mux.ServeHTTP(w, r)

			if w.Code != tc.expectedStatusCode {
				t.Fatalf("Expected status code: %d, got: %d\n", tc.expectedStatusCode, w.Code)
			}

			resMsg := strings.TrimSpace(w.Body.String())
			if tc.expectedMsg != resMsg {
				t.Fatalf("Expected msg: %s, got: %s\n", tc.expectedMsg, resMsg)
			}

			/*------------check to verify document was deleted-----------*/

			if w.Code == http.StatusOK && resMsg == "success" {
				res := shippingAddrColl.FindOne(
					context.Background(),
					bson.M{"_id": id},
				)

				if res.Err() != mongo.ErrNoDocuments {
					t.Fatal("Documented not deleted")
				}
			}
		})
	}
}

func TestHandlePurchaseHistory(t *testing.T) {
	db.Init()
	defer db.Close()
	sessionManager := api.GetSessionManager()
	/*-----------Fake logged in user 1-------------*/

	session1, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "testtest-test-test-test-testtesttest",
	})
	if err != nil {
		t.Fatal("Failed to create session for testuser1")
	}
	userID1, _ := primitive.ObjectIDFromHex("65b094f4a2cb3bf5e40d42d7")
	session1.Store["userid"] = userID1

	/*----------------------tests--------------------*/

	tests := []struct {
		name               string
		sessionID          string
		expectedStatusCode int
	}{
		{
			name:               "Test 1",
			sessionID:          session1.SessionID,
			expectedStatusCode: http.StatusOK,
		},
	}

	server := api.NewServer(":3000")
	server.Use("GET /account/purchase_history", server.HandlePurchaseHistory, api.AuthMiddleware)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/account/purchase_history", nil)
			r.AddCookie(&http.Cookie{
				Name:  api.SESSIONID_COOKIE_NAME,
				Value: tc.sessionID,
			})

			w := httptest.NewRecorder()

			server.Mux.ServeHTTP(w, r)

			if w.Code != tc.expectedStatusCode {
				t.Fatalf("Expected status code: %d, got: %d\n", tc.expectedStatusCode, w.Code)
			}

			if reflect.TypeOf(w.Body.String()).Kind() != reflect.String {
				t.Fatal("Expected return msg to be a string")
			}
		})
	}
}

func TestHandleGETUserFurnitureListings(t *testing.T) {
	db.Init()
	defer db.Close()
	sessionManager := api.GetSessionManager()

	/*-----------Fake logged in user 1-------------*/

	session1, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "testtest-test-test-test-testtesttest",
	})
	if err != nil {
		t.Fatal("Failed to create session for testuser1")
	}
	userID1, _ := primitive.ObjectIDFromHex("65b094f4a2cb3bf5e40d42d7")
	session1.Store["userid"] = userID1

	/*----------------------tests--------------------*/

	tests := []struct {
		name               string
		sessionID          string
		userID             primitive.ObjectID
		expectedStatusCode int
	}{
		{
			name:               "Test 1",
			sessionID:          session1.SessionID,
			userID:             userID1,
			expectedStatusCode: http.StatusOK,
		},
	}

	server := api.NewServer(":3000")
	server.Use("GET /account/furniture_listings", server.HandleGETUserFurnitureListings, api.AuthMiddleware)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/account/furniture_listings", nil)
			r.AddCookie(&http.Cookie{
				Name:  api.SESSIONID_COOKIE_NAME,
				Value: tc.sessionID,
			})

			w := httptest.NewRecorder()

			server.Mux.ServeHTTP(w, r)

			if w.Code != tc.expectedStatusCode {
				t.Fatalf("Expected code: %d, got: %d\n", tc.expectedStatusCode, w.Code)
			}

			var userFurnitureListings []types.FurnitureListing
			err := json.Unmarshal(w.Body.Bytes(), &userFurnitureListings)
			if err != nil {
				t.Fatalf("Error decoding resulting string: %s\n", err.Error())
			}

			// check to make sure each returned listing belongs to the user
			for _, listing := range userFurnitureListings {
				if listing.UserID != tc.userID {
					t.Fatalf("Expected listing of UserID: %s, got: %s\n",
						tc.userID.Hex(), listing.UserID.Hex())
				}
			}
		})
	}

}
