package db

import (
	"backend/types"
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString string = "mongodb://localhost:27017"
const DATABASE_NAME string = "AntiqueFurnitureProject"
const DATABASE_CONTEXT_TIMEOUT time.Duration = 10 * time.Second

var (
	once     sync.Once
	dbClient *mongo.Client
)

// connect with mongoDB
func Init() (*mongo.Client, error) {
	var err error = nil

	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), DATABASE_CONTEXT_TIMEOUT)
		defer cancel()

		clientOptions := options.Client().ApplyURI(connectionString)
		dbClient, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			panic(err)
		}

		log.Println("\x1b[34mConnected to database.\x1b[0m")
	})
	return dbClient, err
}

func GetCollection(collection string) *mongo.Collection {
	return dbClient.Database(DATABASE_NAME).Collection(collection)
}

func Close() error {
	if dbClient != nil {
		err := dbClient.Disconnect(context.Background())
		dbClient = nil
		return err
	}
	return nil
}

/*
This function check whether the specified fieldName in the "users" collection is unique
*/
func CheckFieldUniqueness(fieldName, val string) bool {
	collection := GetCollection("users")
	var result bson.M
	err := collection.FindOne(context.Background(), bson.M{fieldName: val}).Decode(&result)

	if err != nil && err == mongo.ErrNoDocuments {
		return true
	} else {
		return false
	}
}

func InsertIntoUsersCollection(signupInfo types.User) (*mongo.InsertOneResult, error) {
	collection := GetCollection("users")
	return collection.InsertOne(context.Background(), signupInfo)
}

func FindInUsersCollection(loginInfo types.User) *mongo.SingleResult {
	collection := GetCollection("users")
	filter := bson.M{
		"username": loginInfo.Username,
		"password": loginInfo.Password,
	}
	return collection.FindOne(context.Background(), filter)
}

func FindByUsernameInCollection(username string) *mongo.SingleResult {
	collection := GetCollection("users")
	filter := bson.M{
		"username": username,
	}
	return collection.FindOne(context.Background(), filter)
}

func InsertIntoListingsCollection(furnitureListInfo types.FurnitureListing) (*mongo.InsertOneResult, error) {
	collection := GetCollection("listings")
	return collection.InsertOne(context.Background(), furnitureListInfo)
}

/*
This function takes the hex string ID and converts it into an ObjectID so that
it can be used to query the mongoDB to search for the associated listing
*/
func FindByIDInListingsCollection(listingId string) (*mongo.SingleResult, error) {
	collection := GetCollection("listings")
	objID, err := primitive.ObjectIDFromHex(listingId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": objID,
	}
	result := collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return result, nil
}
