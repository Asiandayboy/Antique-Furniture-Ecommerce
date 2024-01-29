package types

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
Type used to save into users collections
*/
type User struct {
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"`
	Email     string             `bson:"email" json:"email"`
	SessionID string             `bson:"sessionid" json:"sessionid"`
	Phone     string             `bson:"phone" json:"phone"`
	UserID    primitive.ObjectID `bson:"_id"`
}

/*
A user can have multiple addresses and can choose
to set a default address to use when buying furniture
*/
type Address struct {
	AddressID primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"userid"`
	State     string             `bson:"state" json:"state"`
	City      string             `bson:"city" json:"city"`
	Street    string             `bson:"street" json:"street"`
	ZipCode   string             `bson:"zipCode" json:"zipCode"`
	Default   bool               `bson:"default" json:"default"`
}

// todo: need to define inputs and outputs for /account
