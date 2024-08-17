package main

import "time"

//go:generate go run gen.go -source $GOFILE -structs Person
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

func NewName(name string) Name {
	return Name(name)
}

func NewAge(age int) Age {
	return Age(age)
}

func NewNumber(number int) Number {
	return Number(number)
}
