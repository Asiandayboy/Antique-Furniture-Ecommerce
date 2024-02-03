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
	Title       string  `bson:"title" json:"title"`
	Description string  `bson:"description" json:"description"`
	Cost        float64 `bson:"cost" json:"cost"`

	// Bed, Table, Desk, Chair, Chest, Nightstand, Cabinet
	Type string `bson:"type" json:"type"`

	// Victorian, English, Baroque, Federal, Rococo, Sheraton, Unknown
	Style string `bson:"style" json:"style"`

	// Mint, Excellent, Good, Worn, Restored, Original Finish
	Condition string `bson:"condition" json:"condition"`

	// Tiger Maple, Cherry, Oak, Walnut, Mahogany, Maple, Chestnut, Paine, Rosewood, Birch
	Material string `bson:"material" json:"material"`

	// images are stored as based64 encoded strings
	Images []string `bson:"images" json:"images"`

	// UserID of the client who created the listing; the owner of the post; the seller
	UserID primitive.ObjectID `bson:"userid"`

	ListingID primitive.ObjectID `bson:"_id"`
}
