package tests

import (
	"backend/db"
	"testing"
)

/*
I have a sample user entry in the database for testing

entry1:
username: "bob"
password: "bob"
email: "bob@gmail.com"

entry2:
username: "johnsmith"
password: "password123"
email: "johnsmith@gmail.com"
*/
func TestCheckFieldUniqueness(t *testing.T) {
	db.Init()
	tests := []struct {
		name     string
		field    string
		payload  string
		expected bool
	}{
		{
			name:     "Test 1",
			field:    "username",
			payload:  "bob",
			expected: false,
		},
		{
			name:     "Test 2",
			field:    "username",
			payload:  "bob1",
			expected: true,
		},
		{
			name:     "Test 3",
			field:    "username",
			payload:  "doglover",
			expected: true,
		},
		{
			name:     "Test 4",
			field:    "username",
			payload:  "johnsmith",
			expected: false,
		},
		{
			name:     "Test 5",
			field:    "username",
			payload:  "john",
			expected: true,
		},
		{
			name:     "Test 6",
			field:    "email",
			payload:  "johnsmith@gmail.com",
			expected: false,
		},
		{
			name:     "Test 7",
			field:    "email",
			payload:  "johnsmith123@gmail.com",
			expected: true,
		},
		{
			name:     "Test 8",
			field:    "email",
			payload:  "bob@gmail.com",
			expected: false,
		},
		{
			name:     "Test 9",
			field:    "email",
			payload:  "bobby@gmail.com",
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			isUnique := db.CheckFieldUniqueness(tc.field, tc.payload)

			if isUnique != tc.expected {
				t.Fatalf("Expected: %v, got: %v\n", tc.expected, isUnique)
			}
		})
	}
	db.Close()
}
