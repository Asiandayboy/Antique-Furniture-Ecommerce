package types

import (
	"backend/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Type used to save into users collections and to
represent client signup and login info, and account info
*/
type User struct {
	UserID    primitive.ObjectID   `bson:"_id"`
	Username  string               `bson:"username" json:"username"`
	Password  string               `bson:"password" json:"password"`
	Email     string               `bson:"email" json:"email"`
	Phone     string               `bson:"phone" json:"phone"`
	SessionID string               `bson:"sessionid"`
	Balance   primitive.Decimal128 `bson:"balance" json:"balance"` // The amount of money from sales in the user's account
}

/*
Provide a negative number to subtract; positive to add
Returns the updated balance
*/
func (u *User) UpdateBalance(amountToAdd float64) float64 {
	currBalance := util.Decimal128ToFloat64(u.Balance)
	newTotal := currBalance + amountToAdd

	dec128 := util.Float64ToDecimal128(newTotal)
	u.Balance = dec128

	return newTotal
}

/*
A user can create multiple shipping addresses and can choose
to set a default address to use when buying furniture

When decoding from JSON into this struct type, provide a custom
unmarshaling to handle the primitive.ObjectID types
*/
type ShippingAddress struct {
	AddressID primitive.ObjectID `bson:"_id,omitempty" json:"addressId"`
	UserID    primitive.ObjectID `bson:"userid" json:"userId"`
	State     string             `bson:"state" json:"state"`
	City      string             `bson:"city" json:"city"`
	Street    string             `bson:"street" json:"street"`
	ZipCode   string             `bson:"zipCode" json:"zipCode"`
	Default   bool               `bson:"default" json:"default"`
}

// todo: need to define inputs and outputs for /account
