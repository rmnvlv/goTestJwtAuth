package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rmnvlv/goTestJwtAuth/routes"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Mail string
	GUID string
}

/* Для проверки поменять dbUrl в routes.db.go */

func main() {
	err := initDB() //Для наполнения бд
	if err != nil {
		log.Fatal(err)
	}

	initRouter()
}

func initRouter() {
	router := http.NewServeMux()

	router.HandleFunc("/get-tokens", routes.GetTokens)
	router.HandleFunc("/refresh-tokens", routes.RefreshToken)
	router.HandleFunc("/protected", routes.Protected) //Тест

	fmt.Println("\nServer starts")

	err := http.ListenAndServe("localhost:80", router)
	if err != nil {
		log.Fatal(err)
	}
}

func initDB() error {
	client, ctx, cancel := routes.ConnectDB()
	defer cancel()

	sessionsCollection := client.Database("testGoRest").Collection("Users")

	mustDelete := bson.D{}
	deleteUsers, err := sessionsCollection.DeleteMany(ctx, mustDelete)
	if err != nil {
		return err
	}
	fmt.Printf("Users delited: %v \n", deleteUsers.DeletedCount)

	user1 := User{"temp@mail", "123od9e0w2kd3o02pw0e2"}
	user2 := User{"temp2@mail", "wew23123sqwew00203ew2"}

	fmt.Printf("Users [mail, guid] for tests: %v \n%v \n", user1, user2)

	_, err = sessionsCollection.InsertOne(ctx, user1)
	if err != nil {
		return err
	}

	_, err = sessionsCollection.InsertOne(ctx, user2)
	if err != nil {
		return err
	}

	return nil
}
