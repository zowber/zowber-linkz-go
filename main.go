package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

type Label struct {
	Id   int
	Name string
}

type Link struct {
	Name         string
	Url          string
	Labels       []Label
	Created_date string
}

func getAllLinks() []Link {

	links := []Link{{
		Name: "First link",
		Url:  "http://example.com/",
		Labels: []Label{
			{
				Id:   1,
				Name: "LabelOne",
			},
			{
				Id:   2,
				Name: "LabelTwo",
			},
		},
		Created_date: "20230909",
	},
		{
			Name: "Second link",
			Url:  "http://example.com/",
			Labels: []Label{
				{
					Id:   1,
					Name: "LabelOne",
				},
				{
					Id:   3,
					Name: "LabelThree",
				},
			},
			Created_date: "20230909",
		}}

	return links
}

var indexHandler = func(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./templates/index.html"))
	temp.Execute(w, getAllLinks())
}

func createLinkHandler(Link) Link {
	log.Println("Creating new link")
	link := Link{
		Name: "First link",
		Url:  "http://example.com/",
		Labels: []Label{
			{
				Id:   1,
				Name: "LabelOne",
			},
			{
				Id:   2,
				Name: "LabelTwo",
			},
		},
		Created_date: "20230909",
	}
	return link
}

func getLinkHandler(linkId int) Link {
	log.Println("Getting link with ID", linkId)
	link := Link{
		Name: "First link",
		Url:  "http://example.com/",
		Labels: []Label{
			{
				Id:   1,
				Name: "LabelOne",
			},
			{
				Id:   2,
				Name: "LabelTwo",
			},
		},
		Created_date: "20230909",
	}
	return link
}

func updateLinkHandler(linkId int, updatedLink Link) Link {
	log.Println("Updating link with ID:", linkId)
	return updatedLink
}

func deleteLinkHandler(linkId int) bool {
	log.Println("Deleting link with ID:", linkId)
	return true
}

func main() {
	log.Println("Sup")

	// API ROUTES
	// /linkz
	// - GET list_all_linkz
	// - POST create_link
	// /links/:linkId
	// - GET read_link
	// - PUT update_link
	// - DELETE delete_link

	/* 	CONTROLLER METHODS
	   	list_all_linkz
	   		uses Find() to get all links
	   		returns all Links */

	http.HandleFunc("/", indexHandler)

	links := getAllLinks()

	for i, l := range links {
		log.Print(i, l)
	}

	/*  create_link
	takes Link
	uses save(Link) to create a new link
	returns the new Link */

	createdLink := Link{
		Name: "Created link",
		Url:  "http://example.com/",
		Labels: []Label{
			{
				Id:   1,
				Name: "LabelOne",
			},
			{
				Id:   2,
				Name: "LabelTwo",
			},
		},
		Created_date: "20230909",
	}

	createLinkHandler(createdLink)

	/*	read_link
		takes Link.Id
		uses findOne(Link.Id) to get a single link
		returns a Link */

	getLinkHandler(1)

	/*	update_link
		takes Link.Id and Link
		uses findOneAndUpdate to modify an existing link
		returns the updated link */

	updatedLink := Link{
		Name: "Updated link",
		Url:  "http://example.com/",
		Labels: []Label{
			{
				Id:   1,
				Name: "LabelOne",
			},
			{
				Id:   2,
				Name: "LabelTwo",
			},
		},
		Created_date: "20230909",
	}

	updateLinkHandler(1, updatedLink)

	/*	delete_link
		takes Link.Id
		uses Remove(Link.Id)
		returns the deleted link?
	*/

	deleteLinkHandler(1)

	http.ListenAndServe(":3000", nil)

}
