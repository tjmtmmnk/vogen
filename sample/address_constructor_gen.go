// Auto-generated constructor for Address
package sample

func NewAddress(number int) *Address {
	return &Address{

		Number: NewAddressNumber(number),
	}
}

func (d AddressNumber) RawValue() int {
	return int(d)
}
