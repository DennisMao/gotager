package tager

type Hello struct {
	Name	string	`json:"Name" json:"Name"`
	Id	int64	`bson:"hi;omitempty" json:"Id" json:"Id"`
}

func HelloPrint() string {
	return ""
}
