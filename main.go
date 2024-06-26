package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"chatty/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbUri = "mongodb://localhost:27017"

func main() {
	client, err := getMongoDbClient(dbUri)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	pingMongoDbClient(client)

	store := NewMongoDBStore(client)
	executeSomeDbCommands(store)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := store.GetAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, user := range users {
			fmt.Fprintln(w, user.Username)
		}
	})

	http.ListenAndServe(":5000", nil)
}

func executeSomeDbCommands(store Datastore) {
	users, err := store.GetAll()
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user.Username)
	}

	randId, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	newUser := User{Username: "user" + fmt.Sprintf("%d", randId)}

	if err = store.Create(newUser); err != nil {
		panic(err)
	}

	users, err = store.GetAll()
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user.Username)
	}

	myUser, err := store.GetByID(newUser.Username)
	if err != nil {
		panic(err)
	}

	fmt.Println(myUser.Username)
}

func getMongoDbClient(uri string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func pingMongoDbClient(client *mongo.Client) {
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Successfully pinged MongoDB!")
}
