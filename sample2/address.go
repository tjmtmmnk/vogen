package sample2

import "testing"

//go:generate go run github.com/tjmtmmnk/vogen -path $GOFILE -structs Address -prefix New -factory true
type Address struct {
	Temp int
}

func BuildAddressTemp(t *testing.T) int {
	return 1
}