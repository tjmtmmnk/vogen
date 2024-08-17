// Auto-generated constructor for Address
package sample

func NewAddress(number int, city string, country string) (*Address, error) {

	t0, err := NewAddressNumber(number)
	if err != nil {
		return nil, err
	}

	t1, err := NewAddressCity(city)
	if err != nil {
		return nil, err
	}

	t2, err := NewAddressCountry(country)
	if err != nil {
		return nil, err
	}

	return &Address{

		Number: t0,

		City: t1,

		Country: t2,
	}, nil
}

func (d AddressNumber) RawValue() int {
	return int(d)
}

func (d AddressCity) RawValue() string {
	return string(d)
}

func (d AddressCountry) RawValue() string {
	return string(d)
}
