package sample

import (
	"testing"
	"time"
)

type (
	PersonName        string
	PersonAge         int
	PersonCatchphrase *string
)

//go:generate go run github.com/tjmtmmnk/vogen -path $GOFILE -structs Person -prefix Parse -factory true
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

func BuildPersonName(t *testing.T) PersonName {
	return PersonName("aa")
}

func BuildPersonAge(t *testing.T) PersonAge {
	return PersonAge(30)
}

func BuildPersonCatchPhrase(t *testing.T) PersonCatchphrase {
	catchPhrase := "Hello, World!"
	return &catchPhrase
}

func BuildPersonCreatedAt(t *testing.T) time.Time {
	return time.Now()
}
