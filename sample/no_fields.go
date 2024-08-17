package sample

//go:generate go run github.com/tjmtmmnk/vogen -source $GOFILE -structs NoFields
type NoFields struct{}
