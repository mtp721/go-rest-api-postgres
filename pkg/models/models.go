package models

type Employee struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	BirthDay  string `json:"birthDay"`
	Gender    string `json:"gender"`
}

type Event struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Date string `json:"date"`
}
