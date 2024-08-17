// Auto-generated constructor for Person
package sample

import "time"

func NewPerson(name string, age int, createdAt time.Time) (*Person, error) {

	t0 := NewPersonName(name)

	t1 := NewPersonAge(age)

	return &Person{

		Name: t0,

		Age: t1,

		CreatedAt: createdAt,
	}, nil
}

func (d PersonName) RawValue() string {
	return string(d)
}

func (d PersonAge) RawValue() int {
	return int(d)
}
