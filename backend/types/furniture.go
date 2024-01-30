package types

import "go.mongodb.org/mongo-driver/bson/primitive"

const Unknown string = "Unknown"

const (
	Bed        string = "Bed"
	Table      string = "Table"
	Desk       string = "Desk"
	Chair      string = "Chair"
	Chest      string = "Chest"
	Nightstand string = "Nightstand"
	Cabinet    string = "Cabinet"
)

const (
	Mint           string = "Mint"
	Excellent      string = "Excellent"
	Good           string = "Good"
	Worn           string = "Worn"
	Restored       string = "Restored"
	OriginalFinish string = "Original Finish"
)

const (
	Victorian string = "Victorian"
	English   string = "English"
	Baroque   string = "Baroque"
	Federal   string = "Federal"
	Rococo    string = "Rococo"
	Sheraton  string = "Sheraton"
)

const (
	TigerMaple string = "Tiger Maple"
	Cherry     string = "Cherry"
	Oak        string = "Oak"
	Walnut     string = "Walnut"
	Mahogany   string = "Mahogany"
	Maple      string = "Maple"
	Chestnut   string = "Chestnut"
	Pine       string = "Pine"
	Rosewood   string = "Rosewood"
	Birch      string = "Birch"
)

/*
This type represents a furniture listing with all
appropriate form details about the furniture
*/
type FurnitureListing struct {
	Title       string  `json:"title"`
	Description string  `json:"desc"`
	Cost        float64 `json:"cost"`

	// Bed, Table, Desk, Chair, Chest, Nightstand, Cabinet
	Type string `json:"type"`

	// Victorian, English, Baroque, Federal, Rococo, Sheraton, Unknown
	Style string `json:"style"`

	// Mint, Excellent, Good, Worn, Restored, Original Finish
	Condition string `json:"condition"`

	// Tiger Maple, Cherry, Oak, Walnut, Mahogany, Maple, Chestnut, Paine, Rosewood, Birch
	Material string `json:"material"`

	// images are stored as based64 encoded strings
	Images []string `json:"images"`

	// UserID of the client who created the listing; the owner of the post
	UserID primitive.ObjectID `bson:"userid"`

	ListingID primitive.ObjectID `bson:"_id"`
}
