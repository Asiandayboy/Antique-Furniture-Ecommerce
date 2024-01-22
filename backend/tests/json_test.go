package tests

import (
	"backend/types"
	"backend/util"
	"encoding/json"
	"io"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type PayloadType int

const (
	UserPayload PayloadType = iota
	FurniturePayload
)

type testUserStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
Testing ReadJSONReq to make sure it is actually decoding
the request body into json and that the json body matches
*/
func TestReadJSONReq(t *testing.T) {
	// formating payloads to simulate a JS obj that has been JSON.stringify()
	payload1Struct := testUserStruct{
		Username: "johnsmith",
		Password: "password123",
	}

	payload2Struct := types.FurnitureListing{
		Title:       "English Tiger maple queen bed",
		Description: "My favorite bed",
		Type:        "Bed",
		Cost:        7500,
		Style:       "English",
		Condition:   "Great",
		Material:    "Tiger Maple",
	}

	payload1Data, err := json.Marshal(payload1Struct)
	if err != nil {
		t.Fatal("Failed to JSON encode payload1struct")
	}
	payload1 := string(payload1Data)

	payload2Data, err := json.Marshal(payload2Struct)
	if err != nil {
		t.Fatal("Failed to JSON encode payload2struct")
	}
	payload2 := string(payload2Data)

	tests := []struct {
		name            string
		payload         string
		expectedPayload interface{}
		expectedError   error
		payloadType     PayloadType
	}{
		{
			name:            "Test 1",
			payload:         payload1,
			expectedPayload: payload1Struct,
			expectedError:   nil,
			payloadType:     UserPayload,
		},
		{
			name:            "Test 2",
			payload:         payload2,
			expectedPayload: payload2Struct,
			expectedError:   nil,
			payloadType:     FurniturePayload,
		},
		{
			name:            "Test 3",
			payload:         `{"username":"bob}`,
			expectedPayload: testUserStruct{}, // expecting an empty struct bc the json is not formatted correctly
			expectedError:   io.ErrUnexpectedEOF,
			payloadType:     UserPayload,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/example", strings.NewReader(tc.payload))

			var decodedUser interface{}
			switch tc.payloadType {
			case UserPayload:
				var user testUserStruct
				err := util.ReadJSONReq[testUserStruct](req, &user)
				if err != tc.expectedError {
					t.Fatal("Error reading JSON request:", err)
				}
				decodedUser = user
			case FurniturePayload:
				var listing types.FurnitureListing
				err := util.ReadJSONReq[types.FurnitureListing](req, &listing)
				if err != tc.expectedError {
					t.Fatal("Error reading JSON request:", err)
				}
				decodedUser = listing
			}

			if !reflect.DeepEqual(decodedUser, tc.expectedPayload) {
				t.Fatalf("Decoded payload does not match expected payload. Got %+v, expected %+v", decodedUser, tc.expectedPayload)
			}

		})
	}
}
