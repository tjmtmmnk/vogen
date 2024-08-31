package sample

import "time"

type (
	PersonName        string
	PersonAge         int
	PersonCatchphrase *string
)

//go:generate go run github.com/tjmtmmnk/vogen -source $GOFILE -structs Person -prefix Parse
type Person struct {
	Name        PersonName
	Age         PersonAge
	CatchPhrase PersonCatchphrase
	CreatedAt   time.Time
}

type NotGenerated struct{}

func ParsePersonName(name string) PersonName {
	return PersonName(name)
}

func ParsePersonAge(age int) PersonAge {
	return PersonAge(age)
}

func ParsePersonCatchphrase(catchPhrase *string) PersonCatchphrase {
	return catchPhrase
}
