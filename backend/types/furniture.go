package types

import "go.mongodb.org/mongo-driver/bson/primitive"

const Unknown string = "Unknown"

// Bed, Table, Desk, Chair, Chest, Nightstand, Cabinet
type FurnitureType = string

// Mint, Excellent, Good, Worn, Restored, Original Finish
type FurnitureCondition = string

// Victorian, English, Baroque, Federal, Rococo, Sheraton, Unknown
type FurnitureStyle = string

// Tiger Maple, Cherry, Oak, Walnut, Mahogany, Maple, Chestnut, Paine, Rosewood, Birch
type FurnitureMaterial = string

const (
	Bed        FurnitureType = "Bed"
	Table      FurnitureType = "Table"
	Desk       FurnitureType = "Desk"
	Chair      FurnitureType = "Chair"
	Chest      FurnitureType = "Chest"
	Nightstand FurnitureType = "Nightstand"
	Cabinet    FurnitureType = "Cabinet"
)

const (
	Mint           FurnitureCondition = "Mint"
	Excellent      FurnitureCondition = "Excellent"
	Good           FurnitureCondition = "Good"
	Worn           FurnitureCondition = "Worn"
	Restored       FurnitureCondition = "Restored"
	OriginalFinish FurnitureCondition = "Original Finish"
)

const (
	Victorian FurnitureStyle = "Victorian"
	English   FurnitureStyle = "English"
	Baroque   FurnitureStyle = "Baroque"
	Federal   FurnitureStyle = "Federal"
	Rococo    FurnitureStyle = "Rococo"
	Sheraton  FurnitureStyle = "Sheraton"
)

const (
	TigerMaple FurnitureMaterial = "Tiger Maple"
	Cherry     FurnitureMaterial = "Cherry"
	Oak        FurnitureMaterial = "Oak"
	Walnut     FurnitureMaterial = "Walnut"
	Mahogany   FurnitureMaterial = "Mahogany"
	Maple      FurnitureMaterial = "Maple"
	Chestnut   FurnitureMaterial = "Chestnut"
	Pine       FurnitureMaterial = "Pine"
	Rosewood   FurnitureMaterial = "Rosewood"
	Birch      FurnitureMaterial = "Birch"
)

/*
This type represents a furniture listing with all
appropriate form details about the furniture
*/
type FurnitureListing struct {
	ListingID   primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Cost        float64            `bson:"cost" json:"cost"`
	Type        FurnitureType      `bson:"type" json:"type"`
	Style       FurnitureStyle     `bson:"style" json:"style"`
	Condition   FurnitureCondition `bson:"condition" json:"condition"`
	Material    FurnitureMaterial  `bson:"material" json:"material"`
	Images      [][]byte           `bson:"images" json:"images"` // images are stored as an array of byte slices
	UserID      primitive.ObjectID `bson:"userid"`               // UserID of the client who created the listing; the owner of the post; the seller
	Bought      bool               `bson:"bought" json:"bought"` // this filed will be used to not render the items that have already been bought
}
