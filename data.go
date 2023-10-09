package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBClient struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func newPortfolioDbClient() (*MongoDBClient, error) {

	connectionStr := DotEnv("DB_URI")
	dbName := "portfolioitems"
	collectionName := "casestudies"

	ctx := context.Background()
	logLvl := options.LogLevel(5)
	loggerOpts := options.Logger().SetComponentLevel(options.LogComponentAll, logLvl)
	clientOpts := options.
		Client().
		ApplyURI(connectionStr).
		SetLoggerOptions(loggerOpts)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Print(err)
		// TODO: Recover from this error
		panic(err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &MongoDBClient{client, collection}, err
}
