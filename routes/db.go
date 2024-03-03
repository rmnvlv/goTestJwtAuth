package routes

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DBURL = "mongodb://localhost:27017"

func UpdateSession(guid, token string) error {
	client, ctx, cancel := ConnectDB()
	defer cancel()

	sessionsCollection := client.Database("testGoRest").Collection("Sessions")
	filter := bson.D{{"guid", guid}}
	update := bson.D{
		{"$set", bson.D{{
			"refreshtoken", token},
		}},
	}

	_, err := sessionsCollection.UpdateOne(ctx, filter, update)
	// fmt.Printf("New session: %v", updatedSession)

	if err != nil {
		return err
	}

	return nil
}

func InsertInDB(guid, token string) error {
	client, ctx, cancel := ConnectDB()
	defer cancel()

	sessionsCollection := client.Database("testGoRest").Collection("Sessions")
	session := Session{guid, token}
	mustDelete := bson.D{{"guid", guid}}

	deleteSession, err := sessionsCollection.DeleteMany(ctx, mustDelete)
	if err != nil {
		return err
	}
	fmt.Printf("Sessions delited: %v", deleteSession.DeletedCount)

	_, err = sessionsCollection.InsertOne(ctx, session)
	if err != nil {
		return err
	}

	fmt.Printf("\nCreated session for: %v\n")

	return nil
}

func SelectSessionFromDbByGuid(guid string) Session {
	client, ctx, cancel := ConnectDB()
	defer cancel()

	usersCollection := client.Database("testGoRest").Collection("Sessions")

	filter := bson.D{{"guid", guid}}
	var user Session
	err := usersCollection.FindOne(ctx, filter).Decode(&user)

	fmt.Println(err)

	return user
}

func SelectUserFromDbByGuid(guid string) Session {
	client, ctx, cancel := ConnectDB()
	defer cancel()

	usersCollection := client.Database("testGoRest").Collection("Users")

	filter := bson.D{{"guid", guid}}
	var user Session
	err := usersCollection.FindOne(ctx, filter).Decode(&user)

	fmt.Println(err)

	return user
}

func SelectFromDbByToken(token string) Session {
	client, ctx, cancel := ConnectDB()
	defer cancel()

	usersCollection := client.Database("testGoRest").Collection("Users")

	filter := bson.D{{"refreshToken", token}}

	user := Session{}
	usersCollection.FindOne(ctx, filter).Decode(&user)

	fmt.Println(usersCollection)

	return user
}

func ConnectDB() (*mongo.Client, context.Context, context.CancelFunc) {

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(DBURL))
	if err != nil {
		fmt.Println("DB connect error: ", err)
	}

	return client, ctx, cancel
}
