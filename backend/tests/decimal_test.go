package tests

import (
	"backend/util"
	"testing"

	// "github.com/stripe/stripe-go/v76"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDecimal128ToFloat64(t *testing.T) {
	payload1, _ := primitive.ParseDecimal128("128.57")
	payload2, _ := primitive.ParseDecimal128("7401.80")

	tests := []struct {
		name     string
		payload  primitive.Decimal128
		expected float64
	}{
		{
			name:     "Test 1",
			payload:  payload1,
			expected: 128.57,
		},
		{
			name:     "Test 2",
			payload:  payload2,
			expected: 7401.80,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := util.Decimal128ToFloat64(tc.payload)

			if res != tc.expected {
				t.Fatalf("Expected: %f, got: %f\n", tc.expected, res)
			}
		})
	}
}

func TestFloat64ToDecimal128(t *testing.T) {
	tests := []struct {
		name     string
		payload  float64
		expected string
	}{
		{
			name:     "Test 1",
			payload:  10500.45,
			expected: "10500.45",
		},
		{
			name:     "Test 2",
			payload:  53025.99,
			expected: "53025.99",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := util.Float64ToDecimal128(tc.payload)
			check, _ := primitive.ParseDecimal128(tc.expected)

			if res != check {
				t.Fatalf("Expected: %s, got: %s\n", tc.expected, res.String())
			}
		})
	}
}
