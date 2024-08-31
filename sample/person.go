package sample

import "time"

type (
	PersonName        string
	PersonAge         int
	PersonCatchphrase *string
)

//go:generate go run github.com/tjmtmmnk/vogen -source $GOFILE -structs Person -prefix Parse -dir sample
type Person struct {
	Name        PersonName
	Age         PersonAge
	CatchPhrase PersonCatchphrase
	CreatedAt   time.Time
	Temp        Temp
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

func ParseTemp(temp int) Temp {
	return Temp(temp)
}
