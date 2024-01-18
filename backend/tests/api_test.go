package tests

import (
	"backend/api"
	"backend/db"
	"backend/types"
	"backend/util"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func trimSpaceAndNewline(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "\n", ""))
}

func TestHandleSignup(t *testing.T) {
	db.Init()
	defer db.Close()
	server := api.NewServer(":3000")
	go server.Start()

	tests := []struct {
		name               string
		method             string
		payload            string
		expectedResMsg     string
		expectedStatusCode int
	}{
		{ // test valid signup; if this test fails, it's bc the payload is already in the db; just delete it and test again
			name:               "Test 1",
			method:             "POST",
			payload:            `{"username": "testuser1", "password": "testpassword1", "email": "test@gmail.com"}`,
			expectedResMsg:     "success",
			expectedStatusCode: http.StatusOK,
		},
		{ // test invalid method
			name:               "Test 2",
			method:             "GET",
			payload:            `{"username": "testuser1", "password": "testpassword1", "email": "test@gmail.com"}`,
			expectedResMsg:     "Request must be a POST request",
			expectedStatusCode: http.StatusBadRequest,
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
			req, err := http.NewRequest(tc.method, "/signup", bytes.NewBufferString(tc.payload))
			if err != nil {
				t.Fatal(err)
			}

			// response recorder response the response
			recorder := httptest.NewRecorder()

			// hit endpoint
			server.HandleSignup(recorder, req)

			if status := recorder.Code; status != tc.expectedStatusCode {
				t.Errorf("Expected code: %v, got: %v", tc.expectedStatusCode, status)
			}
			if res := recorder.Body; trimSpaceAndNewline(res.String()) != tc.expectedResMsg {
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
	go server.Start()

	tests := []struct {
		name               string
		method             string
		payload            string
		expectedResMsg     string
		expectedStatusCode int
	}{
		{ // test valid login
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
			expectedStatusCode: http.StatusBadRequest,
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
			req, err := http.NewRequest(tc.method, "/login", bytes.NewBufferString(tc.payload))
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			server.HandleLogin(recorder, req)

			if status := recorder.Code; status != tc.expectedStatusCode {
				t.Errorf("Expected code: %v, got: %v", tc.expectedStatusCode, status)
			}
			if res := recorder.Body; trimSpaceAndNewline(res.String()) != tc.expectedResMsg {
				t.Errorf("Expected ResMsg: %v, got: %v", tc.expectedResMsg, res.String())
			}
		})
	}
}

/*
This test is to see if the server is setting the cookie
for the session in the response header

This test should be run after TestHandleSignup is called, to ensure
the tested user is added first
*/
func TestHandleLoginCookie(t *testing.T) {
	db.Init()
	defer db.Close()
	server := api.NewServer(":3000")
	go server.Start()

	payload := `{"username": "testuser1", "password": "testpassword1"}`
	req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	server.HandleLogin(recorder, req)

	cookie := recorder.Header().Get("Set-Cookie")
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
  - check this by string comparing the expected body and the mongo document
    and ensuring they're the same

*From the frontend (React TS), the body is sent as one long string
representing a JSON object
*/
func TestHandleListFurnitureJSON(t *testing.T) {

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
		Type:        "Bed",
		Cost:        7500,
		Style:       "English",
		Condition:   "Great",
		Material:    "Tiger Maple",
		Images:      []string{imgData1, imgData2},
	}

	// simulating JSON.stringify(obj) in TS
	payload1, err := json.Marshal(furnitureListing1)
	if err != nil {
		t.Fatal("Failed to encode payload into JSON")
	}
	payload1String := string(payload1)

	// create test cases
	tests := []struct {
		name               string
		payload            string
		expectedStatusCode int
		expectedBody       string
	}{
		{ // test to check if server is receiving the json payload correctly
			name:               "Test 1",
			payload:            payload1String,
			expectedStatusCode: http.StatusOK,
			expectedBody:       payload1String,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := api.NewServer(":3000")
			go server.Start()

			req, err := http.NewRequest("POST", "/list_furniture", bytes.NewBufferString(payload1String))
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			server.HandleListFurniture(recorder, req)

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatal(err)
			}

			// validate JSON body is received
			if string(body) != tc.expectedBody {
				t.Fatalf("Expected: %s, got: %s", tc.expectedBody, string(body))
			}

			// validate that body was inserted into MongoDB correctly

			if status := recorder.Code; status != tc.expectedStatusCode {
				t.Fatalf("Expected status code: %d, got: %v\n", status, tc.expectedStatusCode)
			}
		})
	}
}

/*
Unit test to test ValidateListFormFields
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

}

func TestHandleCheckout(t *testing.T) {

}
