package sample

type (
	AddressNumber  int
	AddressNumber2 AddressNumber
	AddressCity    string
	AddressCountry string
)

//go:generate go run github.com/tjmtmmnk/vogen -source $GOFILE -structs Address
type Address struct {
	Number  AddressNumber
	Number2 AddressNumber2
	City    AddressCity
	Country AddressCountry
}

func NewAddressNumber(number int) (AddressNumber, error) {
	return AddressNumber(number), nil
}

func NewAddressNumber2(number AddressNumber) (AddressNumber2, error) {
	return AddressNumber2(number), nil
}

func NewAddressCity(city string) (AddressCity, error) {
	return AddressCity(city), nil
}

func NewAddressCountry(country string) (AddressCountry, error) {
	return AddressCountry(country), nil
}
