package types

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
Type used to save into users collections
*/
type User struct {
	Username  string             `json:"username"`
	Password  string             `json:"password"`
	Email     string             `json:"email"`
	SessionId string             `json:"sessionid"`
	UserId    primitive.ObjectID `bson:"_id"`
}
