package sample

type (
	AddressNumber int
)

//go:generate go run github.com/tjmtmmnk/vogen -source $GOFILE -structs Address
type Address struct {
	Number AddressNumber
}

func NewAddressNumber(number int) AddressNumber {
	return AddressNumber(number)
}
