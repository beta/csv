// Copyright (c) 2018 Beta Kuang
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package csv_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/beta/csv"
)

const (
	calendarCSV = `first_name,last_name,age,married,phone
John,Smith,25,true,1234567890
Mary,Jane,23,false,9876543210`
	invalidPhoneCalendarCSV = `first_name,last_name,age,married,phone
John,Smith,25,true,12345`
	invalidAgeCalendarCSV = `first_name,last_name,age,married,phone
John,Smith,0,true,1234567890`
)

type Person struct {
	FirstName string `csv:"first_name"`
	LastName  string `csv:"last_name"`
	Age       int    `csv:"age"`
	Married   bool   `csv:"married"`
	Phone     Phone  `csv:"phone"`
}

type PersonValidatingAge struct {
	FirstName string `csv:"first_name"`
	LastName  string `csv:"last_name"`
	Age       int    `csv:"age,age"`
	Married   bool   `csv:"married"`
	Phone     Phone  `csv:"phone"`
}

type Phone string

func (p *Phone) UnmarshalText(text []byte) error {
	var phone = string(text)
	if len(phone) != 10 {
		return fmt.Errorf("%s is not a valid phone number", phone)
	}
	*p = Phone(phone)
	return nil
}

func validateAge(v interface{}) bool {
	if age, ok := v.(int); ok {
		return age > 0
	}
	return false
}

func TestUnmarshal(t *testing.T) {
	var persons []*Person
	var err = csv.Unmarshal([]byte(calendarCSV), &persons)
	if err != nil {
		t.Error(err)
		return
	}
	printPersons(t, persons)
}

func TestUnmarshalCustomType(t *testing.T) {
	var persons []*Person
	var err = csv.Unmarshal([]byte(invalidPhoneCalendarCSV), &persons)
	if err != nil {
		if err.Error() == "12345 is not a valid phone number" {
			t.Log(err)
			return
		}
		t.Error(err)
	}
}

func TestUnmarshalInvalidType(t *testing.T) {
	var persons []Person
	var err = csv.Unmarshal([]byte(calendarCSV), &persons)
	if err != nil {
		if strings.Contains(err.Error(), "is not a pointer to a struct pointer slice") {
			t.Log(err)
			return
		}
		t.Error(err)
	}
}

func TestUnmarshalWithValidator(t *testing.T) {
	var persons []*PersonValidatingAge
	var err = csv.Unmarshal([]byte(invalidAgeCalendarCSV), &persons,
		csv.Validator("age", validateAge))
	if err != nil {
		if strings.Contains(err.Error(), "invalid value") {
			t.Log(err)
			return
		}
		t.Error(err)
	}
}

func printPersons(t *testing.T, persons []*Person) {
	for i, person := range persons {
		t.Logf("Person #%d: { Name: %s %s, Age: %d, Married: %v, Phone: %s }", i, person.FirstName, person.LastName, person.Age, person.Married, person.Phone)
	}
}
