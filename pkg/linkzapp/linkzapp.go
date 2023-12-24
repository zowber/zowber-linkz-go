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

type AppProps struct {
	Settings Settings
}

type Settings struct {
	UserId      int
	Name        string
	ColorScheme string
}

type User struct {
	Id   int
	Name string
}

// for stats

type TotalLinksByMonth struct {
	Month int
	Total int
}

type TotalLinksByYear struct {
	Year   int
	Total  int
	Months []TotalLinksByMonth
}

type TotalLinksByYearAndMonth struct {
	Years []TotalLinksByYear
}

type PopularLabelsByMonth struct {
	Month  int
	Labels []string
}

type PopularLabelsByYear struct {
	Year   int
	Labels []string
	Months []PopularLabelsByMonth
}

type PopularLabelsByYearAndMonth struct {
	Years []PopularLabelsByYear
}
