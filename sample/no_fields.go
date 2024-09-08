package sample

//go:generate go run github.com/tjmtmmnk/vogen -path $GOFILE -structs NoFields -prefix New
type NoFields struct{}
