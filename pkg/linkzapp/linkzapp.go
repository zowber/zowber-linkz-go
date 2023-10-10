package linkzapp

type Label struct {
	Id   int
	Name string
}

type Link struct {
	Id           int
	Name         string
	Url          string
	Labels       []Label
	Created_date string
}
