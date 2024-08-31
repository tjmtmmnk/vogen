package sample

import "testing"

type (
	AddressNumber   int
	AddressNumber2  AddressNumber
	AddressNumber2p *AddressNumber
	AddressCity     string
	AddressCountry  string
	Temp            int
)

//go:generate go run github.com/tjmtmmnk/vogen -source $GOFILE -structs Address -prefix New -dir sample -factory true
type Address struct {
	Number   AddressNumber
	Number2  AddressNumber2
	Number2p AddressNumber2p
	City     AddressCity
	Country  AddressCountry
}

func NewAddressNumber(number int) (AddressNumber, error) {
	return AddressNumber(number), nil
}

func NewAddressNumber2(number int) (AddressNumber2, error) {
	return AddressNumber2(number), nil
}

func NewAddressNumber2p(number *int) (AddressNumber2p, error) {
	var n2p AddressNumber2p
	if number != nil {
		n, err := NewAddressNumber(*number)
		if err != nil {
			return nil, err
		}
		n2p = &n
	}
	return n2p, nil
}

func NewAddressCity(city string) (AddressCity, error) {
	return AddressCity(city), nil
}

func NewAddressCountry(country string) (AddressCountry, error) {
	return AddressCountry(country), nil
}

func BuildAddressNumber(t *testing.T) AddressNumber {
	return AddressNumber(1)
}

func BuildAddressNumber2(t *testing.T) AddressNumber2 {
	return AddressNumber2(2)
}

func BuildAddressNumber2p(t *testing.T) AddressNumber2p {
	temp := AddressNumber(3)
	return &temp
}

func BuildAddressCity(t *testing.T) AddressCity {
	return AddressCity("city")
}

func BuildAddressCountry(t *testing.T) AddressCountry {
	return AddressCountry("country")
}
