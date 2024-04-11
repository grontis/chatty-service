package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBStore struct {
	client *mongo.Client
}

func NewMongoDBStore(client *mongo.Client) *MongoDBStore {
	return &MongoDBStore{client}
}

func (s MongoDBStore) GetAll() ([]*User, error) {
	users := s.client.Database("chatty").Collection("user")
	findOptions := options.Find()
	var results []*User

	cur, err := users.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem User
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return results, err
}

func (s MongoDBStore) GetByID(id string) (*User, error) {
	var result User

	collection := s.client.Database("chatty").Collection("user")
	err := collection.FindOne(context.TODO(), bson.D{{Key: "username", Value: id}}).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}

	return &result, nil
}

func (s MongoDBStore) Create(user User) error {
	users := s.client.Database("chatty").Collection("user")
	_, err := users.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

func (s MongoDBStore) Update(user User) error {
	return nil
}

func (s MongoDBStore) Delete(id string) error {
	return nil
}
