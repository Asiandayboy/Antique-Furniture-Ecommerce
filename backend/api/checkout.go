package api

import (
	"backend/db"
	"backend/types"
	"backend/util"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/stripe/stripe-go/v76"
	stripeSession "github.com/stripe/stripe-go/v76/checkout/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
- Client presses checkout button
- Server reads from request to get data?
- Server create checkout session and redirects
- Server listens for checkout session completed to give receipt
- Server saves appropriate data
*/

const domain string = "http://localhost:5173"

const (
	ErrCheckoutSession = "Error creating checkout session"
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

// This receipt type is what will be stored in the DB
type ReceiptDatabase struct {
	ShippingAddress   ShippingAddress       `bson:"shippingAddress" json:"shippingAddress"`
	PaymentMethod     string                `bson:"paymentMethod" json:"paymentMethod"`
	TotalCost         float32               `bson:"totalCost" json:"totalCost"`
	Items             []ProductItemDatabase `bson:"items" json:"items"` // ID of buyer
	UserID            primitive.ObjectID    `bson:"userid" json:"userId"`
	OrderID           primitive.ObjectID    `bson:"_id" json:"orderId"`
	DatePurchased     time.Time             `bson:"datePurchased" json:"datePurchased"`
	EstimatedDelivery time.Time             `bson:"estimatedDelivery" json:"estimatedDelivery"`
}

// Used with ReceiptResponse
type ProductItemClient struct {
	// ID of the furniture listing
	ListingID string `json:"listingId"`
	// ID of the user who posted the furniture listing; the seller
	SellerID string `json:"sellerId"`
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
	UserID            string              `json:"userId"` // ID of buyer
	OrderID           string              `json:"orderId"`
	DatePurchased     time.Time           `json:"datePurchased"`
	EstimatedDelivery time.Time           `json:"estimatedDelivery"`
}

// func calculateTotalCostTaxed(taxRate float32, prices []float32) float32 {
// 	var sum float32
// 	for _, p := range prices {
// 		sum += p
// 	}
// 	sum += sum * taxRate
// 	return sum
// }

/*
Creates a checkout session from Stripe's API, which redirects the client
to Stripe's hosted page to collect their payment infomation

Stripe servers will then process the payment once the client submits the form
*/
func (s *Server) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	var input CheckoutInfo
	err := util.ReadJSONReq[CheckoutInfo](r, &input)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// query DB with list of listingIDs in shopping cart
	var listingIDsToRetrieve []primitive.ObjectID
	for _, listingID := range input.ShoppingCart {
		objID, err := primitive.ObjectIDFromHex(listingID)
		if err != nil {
			http.Error(w, primitive.ErrInvalidHex.Error(), http.StatusBadRequest)
			return
		}

		listingIDsToRetrieve = append(listingIDsToRetrieve, objID)
	}
	listingsCollection := db.GetCollection("listings")
	filter := bson.M{"_id": bson.M{"$in": listingIDsToRetrieve}}
	cursor, err := listingsCollection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, "Error getting listings", http.StatusBadGateway)
		return
	}

	// extract documents from cursor into an array
	var furnitures []types.FurnitureListing
	for cursor.Next(context.Background()) {
		var furnitureListing types.FurnitureListing
		if err := cursor.Decode(&furnitureListing); err != nil {
			log.Printf("Error decoding document: %v", err)
			continue
		}
		furnitures = append(furnitures, furnitureListing)
	}

	// fmt.Println("Description:", furnitures[0])

	if err := cursor.Err(); err != nil {
		http.Error(w, "Error fetching furnitures", http.StatusBadRequest)
		return
	}

	/*-------------STRIPE-------------*/

	stripe.Key = os.Getenv("STRIPE_TEST_KEY")

	// create Stripe checkout session
	// lineItems is the list of products the customer is buying
	var lineItems []*stripe.CheckoutSessionLineItemParams
	for _, furniture := range furnitures {
		productData := &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
			Name:        &furniture.Title,
			Description: &furniture.Description,
			Metadata: map[string]string{
				"Material":  furniture.Type,
				"Condition": furniture.Material,
				"Style":     furniture.Style,
				"ListingID": furniture.ListingID.Hex(),
			},
		}

		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency:          stripe.String(input.Payment.Currency),
				ProductData:       productData,
				UnitAmountDecimal: &furniture.Cost,
			},
			Quantity: stripe.Int64(1),
		})
	}

	session := r.Context().Value(SessionKey).(*Session)

	params := &stripe.CheckoutSessionParams{
		LineItems:  lineItems,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "/checkout_success"), // frontend page
		CancelURL:  stripe.String(domain + "/checkout_cancel"),  // frontend page

		/*
			This is how we're passing the sessionID; we will access this in the webhook
			handler so we can authenticate the user and fulfill the order by saving
			necessary information associated with the client, like receipts and stuff
		*/
		Metadata: map[string]string{
			"sessionID": session.SessionID,
		},
	}

	checkoutSession, err := stripeSession.New(params)
	if err != nil {
		http.Error(w, ErrCheckoutSession, http.StatusInternalServerError)
		return
	}

	fmt.Println("Checkout session link:", checkoutSession.URL)

	// http.Redirect(w, r, checkoutSession.SuccessURL, http.StatusSeeOther)
}

/*
Processes webhook requests from Stripe's API.

When a stripe checkout is successfully completed, Stripe's server
sends a webhook request to your endpoint, sending data of the event

This handler is used to update the user's account after a successful
checkout and returns necessary information to the client, like their receipt
*/
func (s *Server) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	// extract the Stripe Event obj from the request body
	var event stripe.Event
	if err := util.ReadJSONReq[stripe.Event](r, &event); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse webhook body json: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(event.Type)

	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Error parsing webhook JSON: %v\n", err.Error()),
				http.StatusBadRequest,
			)
		}

		// ERROR: METADATA IS EMPTY :OOO
		fmt.Println(event.Data.Object)

		metadata := checkoutSession.Metadata
		sessionManager := GetSessionManager()
		session, sessionExists := sessionManager.GetSession(metadata["sessionID"])
		fmt.Println(sessionExists, metadata)
		if !sessionExists {
			http.Error(w, "Cannot find sessionID given metadata", http.StatusBadRequest)
			return
		}

		fmt.Printf("Session of user: %s\n", session.SessionID)
	}

	w.WriteHeader(http.StatusOK)
}

// 1. decode request body to get checkout info [DONE]
// 2. fetch documents of shopping cart from listings [DONE]
// 3. calculate total cost with tax [DONE]
// 4. use Stripe to process payment [WIP]
// 5. Transfer money to platform account [DONE]
// 6. delete documents of shopping cart from listings? [WIP]
// 7. generate a receipt for each purchased item [DONE]
// 8. add recept to client's order history [WIP]
// 9. return 200 and receipt [DONE]

// generate the product items
// var productItemsDB []ProductItemDatabase
// for _, furniture := range furitures {
// 	item := ProductItemDatabase{
// 		ListingID: furniture.ListingID,
// 		SellerID:  furniture.UserID,
// 	}
// 	productItemsDB = append(productItemsDB, item)
// }

// generate receipt for client

// add up cost of items in cart to calculate total
// var prices []float32
// for _, furniture := range furitures {
// 	prices = append(prices, float32(furniture.Cost))
// }
// totalCost := calculateTotalCostTaxed(0, prices)

// Start Stripe checkout session

// checkoutSession := &stripe.CheckoutSessionParams{}

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
// var productItemsClient []ProductItemClient
// for _, furniture := range furitures {
// 	item := ProductItemClient{
// 		ListingID: furniture.ListingID.Hex(),
// 		SellerID:  furniture.UserID.Hex(),
// 	}
// 	productItemsClient = append(productItemsClient, item)
// }

// clientReceipt := ReceiptResponse{
// 	ShippingAddress: input.ShippingAddress,
// 	PaymentMethod:   input.Payment.PaymentMethod,
// 	TotalCost:       totalCost,
// 	Items:           productItemsClient,
// 	UserID:          session.Store["userid"].(primitive.ObjectID).Hex(),
// }

// receiptData, err := json.Marshal(clientReceipt)
// if err != nil {
// 	http.Error(w, "Error encoding receipt into JSON", http.StatusInternalServerError)
// 	return
// }

// w.WriteHeader(http.StatusOK)
// w.Write(receiptData)
