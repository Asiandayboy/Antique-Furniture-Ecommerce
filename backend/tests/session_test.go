package tests

import (
	"backend/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSession(t *testing.T) {
	/*
		- create session manager
		- create a session
		- compare memory address of the sessions
	*/
	sessionManager := api.GetSessionManager()

	expectedSession, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "",
	})

	if err != nil {
		t.Fatal("Failed to create new session:", err)
	}

	actualSession, exists := sessionManager.GetSession(expectedSession.SessionId)
	if !exists {
		t.Fatal("Session does not exist:", expectedSession.SessionId)
	}

	// compare memory address of the session
	if expectedSession != actualSession {
		t.Fatalf("Expected memAddr: %p, got: %p\n", expectedSession, actualSession)
	}
}

func TestIsLoggedIn(t *testing.T) {
	sessionManager := api.GetSessionManager()

	// simulating a loggedIn user
	session, err := sessionManager.CreateSession(api.SessionTemplate{
		SessionID: "",
	})

	if err != nil {
		t.Fatal("Fail to generate new sessionID when testing")
	}

	tests := []struct {
		name           string
		payload        string
		expectedReturn bool
	}{
		{
			name:           "Test 1",
			payload:        "1234597123",
			expectedReturn: false,
		},
		{
			name:           "Test 2",
			payload:        session.SessionId,
			expectedReturn: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/foo", nil)
			// add cookie to request
			r.AddCookie(&http.Cookie{
				Name:  api.SESSIONID_COOKIE_NAME,
				Value: tc.payload,
			})

			_, exists := sessionManager.IsLoggedIn(r)
			if exists != tc.expectedReturn {
				t.Fatalf("Expected return: %v, got: %v\n", tc.expectedReturn, exists)
			}

		})
	}
}
