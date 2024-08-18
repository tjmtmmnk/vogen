package sample

import "time"

type (
	PersonName        string
	PersonAge         int
	PersonCatchphrase *string
)

//go:generate go run github.com/tjmtmmnk/vogen -source $GOFILE -structs Person
type Person struct {
	Name        PersonName
	Age         PersonAge
	CatchPhrase PersonCatchphrase
	CreatedAt   time.Time
}

type NotGenerated struct{}

func NewPersonName(name string) PersonName {
	return PersonName(name)
}

func NewPersonAge(age int) PersonAge {
	return PersonAge(age)
}

func NewPersonCatchphrase(catchPhrase *string) PersonCatchphrase {
	return catchPhrase
}
