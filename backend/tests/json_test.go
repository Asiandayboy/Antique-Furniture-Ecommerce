package tests

import (
	"backend/util"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

type testUserStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func compareJSONObj(a, b *testUserStruct) bool {
	jsonA, errA := json.Marshal(a)
	jsonB, errB := json.Marshal(b)
	if errA != nil || errB != nil {
		return false
	}
	return string(jsonA) == string(jsonB)
}

func TestReadJSONReq(t *testing.T) {
	payload := `{
		"username": "johnsmith",
		"password": "password123"
	}`

	// create new HTTP request with sample data
	req, err := http.NewRequest("POST", "/example", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	// simulate HTTP request
	req.Header.Set("Content-Type", "application/json")

	var user testUserStruct
	err = util.ReadJSONReq(req, &user)
	if err != nil {
		t.Fatalf("Error reading JSON request: %v", err)
	}

	expectedUser := testUserStruct{
		Username: "johnsmith",
		Password: "password123",
	}
	// compare
	if !compareJSONObj(&user, &expectedUser) {
		t.Errorf("Expected %v, got %v", expectedUser, user)
	}
}
