package linkzapp

type Label struct {
	Id   *int    `bson:"id,omitempty"`
	Name string `bson:"name,omitempty"`
}

type Link struct {
	Id        *int    `bson:"_id,omitempty"`
	Name      string  `bson:"name,omitempty"`
	Url       string  `bson:"url,omitempty"`
	Labels    []Label `bson:"labels,omitempty"`
	CreatedAt int     `bson:"createdat,omitempty"`
}
