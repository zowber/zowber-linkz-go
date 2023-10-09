package main

var linksStub = []Link{{
	Id:   1,
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

var linkStub = Link{
	Id:   2,
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
