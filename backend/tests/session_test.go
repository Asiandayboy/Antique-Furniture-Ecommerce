package tests

import (
	"backend/api"
	"testing"
)

func TestGetSession(t *testing.T) {
	/*
		- create session manager
		- create a session
		- compare memory address of the sessions
	*/
	sessionManager := api.GetSessionManager()

	expectedSession, err := sessionManager.CreateSession()
	if err != nil {
		t.Fatal("Failed to create new session:", err)
	}

	actualSession, err := sessionManager.GetSession(expectedSession.SessionId.String())
	if err != nil {
		t.Fatal("Failed to create new session:", err)
	}

	// compare memory address of the session
	if expectedSession != actualSession {
		t.Fatalf("Expected memAddr: %p, got: %p\n", expectedSession, actualSession)
	}
}
