package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Database

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ConnectMongoDB() {
	dbUrl := os.Getenv("DB_URL")
	client, err := mongo.NewClient(options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Connect(context.Background()); err != nil {
		log.Fatal(err)
	} else {
		if err = client.Ping(context.Background(), readpref.Primary()); err == nil {
			log.Println("Successfully connected to mongodb")
		} else {
			log.Fatal(err)
		}
	}

	db = client.Database("rentless")
	fmt.Println(db.Name())
	fmt.Println("Connected")
}

func InsertOne(data interface{}) *mongo.InsertOneResult {
	var insertResult *mongo.InsertOneResult
	var err error
	if db != nil {
		collection := db.Collection("product")
		if insertResult, err = collection.InsertOne(context.Background(), data); err != nil {
			log.Fatal(err)
		}
	}
	return insertResult
}

func DeleteOne(filter interface{}) (*mongo.DeleteResult, error) {
	var deleteResult *mongo.DeleteResult
	var err error
	if db != nil {
		collection := db.Collection("product")
		if deleteResult, err = collection.DeleteOne(context.Background(), filter); err != nil {
			return nil, err
		}
	}
	return deleteResult, nil
}

func UpdateOne(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	var updateResult *mongo.UpdateResult
	var err error
	if db != nil {
		collection := db.Collection("product")
		if updateResult, err = collection.UpdateOne(context.Background(), filter, update); err != nil {
			return nil, err
		}
	}
	return updateResult, err
}

func GetAllProduct() (*mongo.Cursor, error) {
	var cursor *mongo.Cursor
	var err error
	if db != nil {
		collection := db.Collection("product")
		if cursor, err = collection.Find(context.Background(), bson.M{}); err != nil {
			return nil, err
		}
	}
	return cursor, nil
}
