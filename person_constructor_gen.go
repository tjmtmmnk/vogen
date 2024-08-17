// Auto-generated constructor for Person
package main

import "time"

func NewPerson(name string, age int, createdAt time.Time) *Person {
	return &Person{

		Name: NewName(name),

		Age: NewAge(age),

		CreatedAt: createdAt,
	}
}
