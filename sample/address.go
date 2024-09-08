package sample

import (
	"testing"
	"time"

	"github.com/tjmtmmnk/vogen/sample2"
)

type (
	AddressNumber   int
	AddressNumber2  AddressNumber
	AddressNumber2p *AddressNumber
	AddressCity     string
	AddressCountry  string
	TempFunc        func() int
	TempSlice       []int
	TempSliceP      []*int
	TempMap         map[string]int
	Temp2           sample2.Temp
	TempTime        time.Time
)

//go:generate go run github.com/tjmtmmnk/vogen -source $GOFILE -structs Address -prefix New -dir sample -factory true
type Address struct {
	Number     AddressNumber
	Number2    AddressNumber2
	Number2p   AddressNumber2p
	City       AddressCity
	Country    AddressCountry
	Temp       sample2.Temp
	TempFunc   TempFunc
	TempSlice  TempSlice
	TempSliceP TempSliceP
	TempMap    TempMap
	Temp2      Temp2
	TempTime   TempTime
}

type Address2 struct {
	Number AddressNumber
}

func NewAddressNumber(number int) (AddressNumber, error) {
	return AddressNumber(number), nil
}

func NewAddressNumber2(number int) (AddressNumber2, error) {
	return AddressNumber2(number), nil
}

func NewAddressCity(city string) (AddressCity, error) {
	return AddressCity(city), nil
}

func NewAddressCountry(country string) (AddressCountry, error) {
	return AddressCountry(country), nil
}

func NewAddressTemp2(temp2 sample2.Temp) Temp2 {
	return Temp2(temp2)
}

func NewAddressTempTime(tempTime time.Time) TempTime {
	return TempTime(tempTime)
}

func NewAddress2Number(number int) (AddressNumber, error) {
	return NewAddressNumber(number)
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

func BuildAddressTemp(t *testing.T) sample2.Temp {
	return sample2.Temp(1)
}

func BuildAddressTempFunc(t *testing.T) TempFunc {
	return func() int {
		return 1
	}
}

func BuildAddressTempSlice(t *testing.T) TempSlice {
	return []int{1, 2, 3}
}

func BuildAddressTempSliceP(t *testing.T) TempSliceP {
	temp := 1
	return []*int{&temp}
}

func BuildAddressTempMap(t *testing.T) TempMap {
	return map[string]int{"key": 1}
}

func BuildAddressTemp2(t *testing.T) Temp2 {
	return Temp2(sample2.Temp(1))
}

func BuildAddressTempTime(t *testing.T) TempTime {
	return TempTime(time.Now())
}

func BuildAddress2Number(t *testing.T) AddressNumber {
	return BuildAddressNumber(t)
}
