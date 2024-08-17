package sample

import "time"

//go:generate go run github.com/tjmtmmnk/vogen -source $GOFILE -structs Person,Address,NoFields
type (
	Name   string
	Age    int
	Number int
)

type Person struct {
	Name      Name
	Age       Age
	CreatedAt time.Time
}

type Address struct {
	Number Number
}

type NoFields struct{}

type NotGenerated struct{}

func NewName(name string) Name {
	return Name(name)
}

func NewAge(age int) Age {
	return Age(age)
}

func NewNumber(number int) Number {
	return Number(number)
}
