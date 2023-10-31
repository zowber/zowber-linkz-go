package mongo

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DotEnv(key string) string {
	if os.Getenv(key) == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}
	return os.Getenv(key)
}

type MongoDBClient struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewPortfolioDbClient() (*MongoDBClient, error) {
	connectionStr := DotEnv("DB_URI")
	dbName := "links"
	collectionName := "links"

	ctx := context.Background()
	logLvl := options.LogLevel(0)
	loggerOpts := options.Logger().SetComponentLevel(options.LogComponentAll, logLvl)
	clientOpts := options.
		Client().
		ApplyURI(connectionStr).
		SetLoggerOptions(loggerOpts)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Print(err)
		panic(err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &MongoDBClient{client, collection}, err
}

func (m *MongoDBClient) All() ([]*linkzapp.Link, error) {
	ctx := context.Background()
	filter := bson.M{}
	opts := options.Find()
	opts.SetSort(bson.M{"createdat": -1})

	var links []*linkzapp.Link
	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return links, err
	}

	if err := cursor.All(ctx, &links); err != nil {
		return links, err
	}

	return links, nil
}

func (m *MongoDBClient) One(linkId primitive.ObjectID) (*linkzapp.Link, error) {
	ctx := context.Background()
	filter := bson.M{"_id": linkId}
	opts := options.FindOne()

	var link *linkzapp.Link
	err := m.collection.FindOne(ctx, filter, opts).Decode(&link)
	if err != nil {
		return &linkzapp.Link{}, err
	}

	return link, nil
}

func (m *MongoDBClient) Insert(link *linkzapp.Link) (*linkzapp.Link, error) {
	ctx := context.Background()
	opts := options.InsertOne()

	result, err := m.collection.InsertOne(ctx, link, opts)
	if err != nil {
		return &linkzapp.Link{}, err
	}

	newLink, err := m.One(result.InsertedID.(primitive.ObjectID))
	if err != nil {
		return &linkzapp.Link{}, err
	}

	return newLink, nil
}

func (m *MongoDBClient) Update(link *linkzapp.Link) (*linkzapp.Link, error) {
	ctx := context.Background()
	filter := bson.M{"_id": link.Id}
	replace := bson.M{"name": link.Name, "url": link.Url, "labels": link.Labels}
	opts := options.FindOneAndReplace().SetReturnDocument(options.After)

	var updatedLink *linkzapp.Link

	err := m.collection.FindOneAndReplace(ctx, filter, replace, opts).Decode(&updatedLink)
	if err != nil {
		return &linkzapp.Link{}, err
	}

	return updatedLink, nil
}

func (m *MongoDBClient) Delete(oid primitive.ObjectID) error {
	ctx := context.Background()
	filter := bson.M{"_id": oid}
	opts := options.FindOneAndDelete()

	var deletedLink *linkzapp.Link

	err := m.collection.FindOneAndDelete(ctx, filter, opts).Decode(&deletedLink)
	if err != nil {
		return err
	}

	return nil
}
