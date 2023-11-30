package linkzapp

type Label struct {
	Id   *int   `json:"id"`   //`bson:"id,omitempty"`
	Name string `json:"name"` //`bson:"name,omitempty"`
}

type Link struct {
	Id        *int    `json:"id"`        //`bson:"_id,omitempty"`
	Name      string  `json:"name"`      //`bson:"name,omitempty"`
	Url       string  `json:"url"`       //`bson:"url,omitempty"`
	Labels    []Label `json:"labels"`    //`bson:"labels,omitempty"`
	CreatedAt int     `json:"createdat"` //`bson:"createdat,omitempty"`
}

type Settings struct {
    UserId int
    Name string
    ColorScheme string
}

type User struct {
    Id int
    Name string
}
