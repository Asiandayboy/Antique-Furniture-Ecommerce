package tests

import (
	"backend/api"
	"backend/db"
	"backend/types"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
This isn't really a proper test. I just wrote it
so I can actually execute the code and see if it's working or not.
I verified that it work by checking if I got the email
*/
func TestSendNewListingNotificationEmail(t *testing.T) {
	tests := []struct {
		name  string
		input types.FurnitureListing
	}{
		{
			name: "Test 1",
			input: types.FurnitureListing{
				ListingID:   primitive.NewObjectID(),
				Title:       "Bobby",
				Description: "Flay",
				Cost:        500,
				Type:        "Chair",
				Style:       "Boring",
				Condition:   "Great",
				Material:    "Oak",
				Bought:      false,
				UserID:      primitive.NewObjectID(),
			},
		},
	}

	db.Init()
	defer db.Close()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := api.SendNewListingNotificationEmail(tc.input)
			if err != nil {
				t.Fatalf("Test: Error occurred while sending emails: %s\n", err.Error())
			}
		})
	}
}
