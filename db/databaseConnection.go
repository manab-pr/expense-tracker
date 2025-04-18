package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {

	// if _, err := os.Stat(".env"); err == nil {
	//     // Load the .env file
	//     if err := godotenv.Load(); err != nil {
	//         log.Fatalf("Error loading .env file: %v", err)
	//     }
	// }

	MongoDb := os.Getenv("MONGO_URI")

	client, err := mongo.NewClient((options.Client().ApplyURI(MongoDb)))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to database....")
	return client

}

var Client *mongo.Client = DBInstance()

func OpenCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("expanse-tracker").Collection(CollectionName)

	return collection
}
