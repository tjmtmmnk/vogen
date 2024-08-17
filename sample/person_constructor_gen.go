// Auto-generated constructor for Person
package sample

import "time"

func NewPerson(name string, age int, catchPhrase PersonCatchphrase, createdAt time.Time) (*Person, error) {

	t0 := NewPersonName(name)

	t1 := NewPersonAge(age)

	return &Person{

		Name: t0,

		Age: t1,

		CatchPhrase: catchPhrase,

		CreatedAt: createdAt,
	}, nil
}

type rawPerson struct {
	Name string

	Age int

	CatchPhrase PersonCatchphrase

	CreatedAt time.Time
}

func (d Person) RawValue() rawPerson {
	return rawPerson{

		Name: d.Name.RawValue(),

		Age: d.Age.RawValue(),

		CatchPhrase: d.CatchPhrase,

		CreatedAt: d.CreatedAt,
	}
}

func (d PersonName) RawValue() string {
	return string(d)
}

func (d PersonAge) RawValue() int {
	return int(d)
}
