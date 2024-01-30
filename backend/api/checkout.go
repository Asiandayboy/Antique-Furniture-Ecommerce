package api

import (
	"backend/db"
	"backend/types"
	"backend/util"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// This struct is to be used with Stripe API
type PaymentInfo struct {
	StripeToken   string  `json:"stripeToken"`
	PaymentMethod string  `json:"paymentMethod"`
	Amount        float32 `json:"amount"`
	Currency      string  `json:"currency"`
}

type ShippingAddress struct {
	State   string `json:"state"`
	City    string `json:"city"`
	Street  string `json:"street"`
	ZipCode string `json:"zipCode"`
}

type CheckoutInfo struct {
	ShoppingCart    []string        `json:"shoppingCart"`
	Payment         PaymentInfo     `json:"paymentInfo"`
	ShippingAddress ShippingAddress `json:"shippingAddress"`
}

// Used with ReceiptDatabase
type ProductItemDatabase struct {
	// ID of the furniture listing
	ListingID primitive.ObjectID `bson:"listingid" json:"listingId"`
	// ID of the user who posted the furniture listing; the seller
	SellerID primitive.ObjectID `bson:"sellerid" json:"sellerId"`
}

// Used with ReceiptResponse
type ProductItemClient struct {
	// ID of the furniture listing
	ListingID string `json:"listingId"`
	// ID of the user who posted the furniture listing; the seller
	SellerID string `json:"sellerId"`
}

// This receipt type is what will be stored in the DB
type ReceiptDatabase struct {
	ShippingAddress   ShippingAddress       `bson:"shippingAddress" json:"shippingAddress"`
	PaymentMethod     string                `bson:"paymentMethod" json:"paymentMethod"`
	TotalCost         float32               `bson:"totalCost" json:"totalCost"`
	Items             []ProductItemDatabase `bson:"items" json:"items"`
	UserID            primitive.ObjectID    `bson:"userid" json:"userId"`
	OrderID           primitive.ObjectID    `bson:"_id" json:"orderId"`
	DatePurchased     time.Time             `bson:"datePurchased" json:"datePurchased"`
	EstimatedDelivery time.Time             `bson:"estimatedDelivery" json:"estimatedDelivery"`
}

/*
This receipt type is sent to the client, where the ObjectIDs are
sent as hexadecimal strings rather than as a primitive.ObjecttID
*/
type ReceiptResponse struct {
	ShippingAddress   ShippingAddress     `json:"shippingAddress"`
	PaymentMethod     string              `json:"paymentMethod"`
	TotalCost         float32             `json:"totalCost"`
	Items             []ProductItemClient `json:"items"`
	UserID            string              `json:"userId"`
	OrderID           string              `json:"orderId"`
	DatePurchased     time.Time           `json:"datePurchased"`
	EstimatedDelivery time.Time           `json:"estimatedDelivery"`
}

func calculateTotalCostTaxed(taxRate float32, prices []float32) float32 {
	var sum float32
	for _, p := range prices {
		sum += p
	}
	sum += sum * taxRate
	return sum
}

/*
Processes checkout requests and handles payment. If the checkout is successful,
the client should receive a 200 code, and a receipt of their order
*/
func (s *Server) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	// 1. decode request body to get checkout info
	// 2. fetch documents of shopping cart from listings
	// 3. calculate total cost with tax
	// 4. use Stripe to process payment
	// 5. delete documents of shopping cart from listings? [or no]
	// 6. generate a receipt for each purchased item
	// 7. add recept to client's order history
	// 8. return 200 and receipt

	session := r.Context().Value(SessionKey).(*Session)

	var input CheckoutInfo
	err := util.ReadJSONReq[CheckoutInfo](r, &input)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// query DB with list of listingIDs in shopping cart
	var listingIDsToRetrieve []primitive.ObjectID
	for _, listingID := range input.ShoppingCart {
		objID, err := primitive.ObjectIDFromHex(string(listingID))
		if err != nil {
			http.Error(w, primitive.ErrInvalidHex.Error(), http.StatusBadRequest)
			return
		}

		listingIDsToRetrieve = append(listingIDsToRetrieve, objID)
	}
	filter := bson.M{"_id": bson.M{"$in": listingIDsToRetrieve}}
	listingsCollection := db.GetCollection("listings")
	cursor, err := listingsCollection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, "Error getting listings", http.StatusBadGateway)
		return
	}

	// add the listings to an array
	var furitures []types.FurnitureListing
	for cursor.Next(context.Background()) {
		var furnitureListing types.FurnitureListing
		if err := cursor.Decode(&furnitureListing); err != nil {
			log.Printf("Error decoding document: %v", err)
			continue
		}
		furitures = append(furitures, furnitureListing)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Error fetching furnitures", http.StatusBadRequest)
		return
	}

	// generate the product items
	// var productItemsDB []ProductItemDatabase
	// for _, furniture := range furitures {
	// 	item := ProductItemDatabase{
	// 		ListingID: furniture.ListingID,
	// 		SellerID:  furniture.UserID,
	// 	}
	// 	productItemsDB = append(productItemsDB, item)
	// }

	var productItemsClient []ProductItemClient
	for _, furniture := range furitures {
		item := ProductItemClient{
			ListingID: furniture.ListingID.Hex(),
			SellerID:  furniture.UserID.Hex(),
		}
		productItemsClient = append(productItemsClient, item)
	}

	// calculate total cost with tax
	var prices []float32
	for _, furniture := range furitures {
		prices = append(prices, float32(furniture.Cost))
	}
	totalCost := calculateTotalCostTaxed(0, prices)

	// genereate receipt for DB
	// dbReceipt := ReceiptDatabase{
	// 	ShippingAddress: input.ShippingAddress,
	// 	PaymentMethod:   input.Payment.PaymentMethod,
	// 	TotalCost:       totalCost,
	// 	Items:           input.ShoppingCart,
	// 	UserID:          session.Store["userid"].(primitive.ObjectID),
	// }

	// save dbReceipt into DB

	// generate and send receipt to client
	clientReceipt := ReceiptResponse{
		ShippingAddress: input.ShippingAddress,
		PaymentMethod:   input.Payment.PaymentMethod,
		TotalCost:       totalCost,
		Items:           productItemsClient,
		UserID:          session.Store["userid"].(primitive.ObjectID).Hex(),
	}

	receiptData, err := json.Marshal(clientReceipt)
	if err != nil {
		http.Error(w, "Error encoding receipt into JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(receiptData)
}
