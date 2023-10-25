package linkzapp

import "go.mongodb.org/mongo-driver/bson/primitive"

type Label struct {
	Id   int    `bson:"id,omitempty"`
	Name string `bson:"name,omitempty"`
}

type Link struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	Url       string             `bson:"url,omitempty"`
	Labels    []Label            `bson:"labels,omitempty"`
	CreatedAt int64              `bson:"createdat,omitempty"`
}
