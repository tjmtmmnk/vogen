package sample

//go:generate go run github.com/tjmtmmnk/vogen -source $GOFILE -structs NoFields -prefix New -dir sample
type NoFields struct{}
