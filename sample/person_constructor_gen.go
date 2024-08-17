// Auto-generated constructor for Person
package sample

import "time"

func NewPerson(name string, age int, createdAt time.Time) *Person {
	return &Person{

		Name: NewPersonName(name),

		Age: NewPersonAge(age),

		CreatedAt: createdAt,
	}
}

func (d PersonName) RawValue() string {
	return string(d)
}

func (d PersonAge) RawValue() int {
	return int(d)
}
