// Copyright (c) 2018 Beta Kuang
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package csv

import (
	"encoding"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

const csvTagName = "csv"

var textUnmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

var (
	// Validator adds a new validator functions for validating a CSV value while
	// unmarshaling a document.
	Validator = func(name string, validator func(interface{}) bool) Setting {
		return func(r *rule) {
			if r.validators == nil {
				r.validators = make(map[string]func(interface{}) bool)
			}
			r.validators[name] = validator
		}
	}
)

// Unmarshal parses a CSV document and stores the result in the struct slice
// pointed to by dest. If dest is nil or not a pointer to a struct slice,
// Unmarshal returns an InvalidUnmarshalError.
func Unmarshal(data []byte, dest interface{}, settings ...Setting) error {
	var v = reflect.ValueOf(dest)
	if v.IsNil() {
		return &InvalidUnmarshalError{Type: nil}
	}
	if v.Type().Kind() != reflect.Ptr || v.Type().Elem().Kind() != reflect.Slice ||
		v.Type().Elem().Elem().Kind() != reflect.Ptr || v.Type().Elem().Elem().Elem().Kind() != reflect.Struct {
		return &InvalidUnmarshalError{Type: reflect.TypeOf(dest)}
	}

	var u = newUnmarshaler(data, dest, settings...)
	return u.unmarshal()
}

func newUnmarshaler(data []byte, dest interface{}, settings ...Setting) *unmarshaler {
	var u = &unmarshaler{
		rule:     defaultRule,
		data:     data,
		dest:     dest,
		settings: settings,
	}
	for _, setting := range settings {
		setting(&u.rule)
	}
	return u
}

type unmarshaler struct {
	rule rule

	data     []byte
	dest     interface{}
	settings []Setting

	fieldMap map[string]*field // Key is the CSV header name of the field.
}

func (u *unmarshaler) prepareFields() {
	// u.dest is a pointer to struct pointer slice.
	var structType = reflect.TypeOf(u.dest).Elem().Elem().Elem()
	var fieldMap = make(map[string]*field, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		var structField = structType.Field(i)
		if tag, exist := structField.Tag.Lookup(csvTagName); exist {
			var tagParts = strings.Split(tag, ",")
			var csvName = tagParts[0]
			if csvName == "-" {
				continue
			}

			var validatorNames = make([]string, 0, len(tagParts)-1)
			for i := 1; i < len(tagParts); i++ {
				validatorNames = append(validatorNames, tagParts[i])
			}

			var field = &field{
				Name:           structField.Name,
				Type:           structField.Type,
				CSVName:        csvName,
				ValidatorNames: validatorNames,
			}
			fieldMap[csvName] = field
		}
	}
	u.fieldMap = fieldMap
}

// Info of a field in the target struct.
type field struct {
	Name           string
	Type           reflect.Type
	CSVName        string
	ValidatorNames []string
}

func (u *unmarshaler) unmarshal() error {
	u.prepareFields()

	var settings = append(u.settings, Header(true))
	var scanner = NewScanner(u.data, settings...)
	header, rows, err := scanner.Scan()
	if err != nil {
		return err
	}

	var sliceV = reflect.ValueOf(u.dest).Elem() // u.dest is a pointer to struct pointer slice.
	for rowIndex, row := range rows {
		var rowCount = rowIndex + 1
		if rowCount > sliceV.Cap() {
			// Grow slice.
			var newCap = sliceV.Cap() + sliceV.Cap()/2
			if newCap < 4 {
				newCap = 4
			}
			var newSliceV = reflect.MakeSlice(sliceV.Type(), sliceV.Len(), newCap)
			reflect.Copy(newSliceV, sliceV)
			sliceV.Set(newSliceV)
		}
		if rowCount >= sliceV.Len() {
			sliceV.SetLen(rowCount)
		}

		var obj = reflect.New(sliceV.Type().Elem().Elem())
		sliceV.Index(rowIndex).Set(obj)
		err = u.unmarshalRecord(sliceV.Index(rowIndex), header, row)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *unmarshaler) unmarshalRecord(dest reflect.Value, header []string, row []string) error {
	for i, value := range row {
		var name = header[i]
		field, exist := u.fieldMap[name]
		if !exist {
			continue
		}

		var err = u.unmarshalField(field, dest.Elem().FieldByName(field.Name), value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *unmarshaler) unmarshalField(field *field, dest reflect.Value, value string) error {
	// Validation.
	for _, validatorName := range field.ValidatorNames {
		validator, exist := u.rule.validators[validatorName]
		if !exist {
			return fmt.Errorf("csv: cannot find validator %s", validatorName)
		}
		if !validator(value) {
			return fmt.Errorf("csv: invalid value %s for field %s", value, field.Name)
		}
	}

	if tu, ok := dest.Interface().(encoding.TextUnmarshaler); ok {
		return tu.UnmarshalText([]byte(value))
		// dest.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value))
	}

	if dest.CanAddr() {
		if tu, ok := dest.Addr().Interface().(encoding.TextUnmarshaler); ok {
			return tu.UnmarshalText([]byte(value))
		}
	}

	var k = dest.Type().Kind()
	if reflect.Int <= k && k <= reflect.Uint64 {
		return u.unmarshalInt(dest, value)
	}
	switch k {
	case reflect.Bool:
		return u.unmarshalBool(dest, value)
	case reflect.Float32, reflect.Float64:
		return u.unmarshalFloat(dest, value)
	case reflect.String:
		return u.unmarshalString(dest, value)
	}
	return fmt.Errorf("csv: unsupported Go type %s", dest.Type().String())
}

func (u *unmarshaler) unmarshalInt(dest reflect.Value, value string) error {
	intVal, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	var uintVal = uint64(intVal)

	// Check integer range.
	var outOfRange = false
	var unsigned = false

	if k := dest.Kind(); k >= reflect.Int && k <= reflect.Int64 {
		switch k {
		case reflect.Int, reflect.Int64:
			// No checking needed.
		case reflect.Int8:
			outOfRange = intVal < math.MinInt8 || intVal > math.MaxInt8
		case reflect.Int16:
			outOfRange = intVal < math.MinInt16 || intVal > math.MaxInt16
		case reflect.Int32:
			outOfRange = intVal < math.MinInt32 || intVal > math.MaxInt32
		}
	} else if k >= reflect.Uint && k <= reflect.Uint64 {
		unsigned = true
		switch k {
		case reflect.Uint, reflect.Uint64:
			// No checking needed.
		case reflect.Int8:
			outOfRange = intVal < 0 || uintVal > math.MaxUint8
		case reflect.Int16:
			outOfRange = intVal < 0 || uintVal > math.MaxUint16
		case reflect.Int32:
			outOfRange = intVal < 0 || uintVal > math.MaxUint32
		}
	}

	if outOfRange {
		return fmt.Errorf("csv: value %s is out of range for type %s", value, dest.Type().String())
	}
	if !unsigned {
		dest.SetInt(intVal)
	} else {
		dest.SetUint(uintVal)
	}
	return nil
}

func (u *unmarshaler) unmarshalBool(dest reflect.Value, value string) error {
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}

	dest.SetBool(boolVal)
	return nil
}

func (u *unmarshaler) unmarshalFloat(dest reflect.Value, value string) error {
	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}

	// Check float range.
	var outOfRange = false
	switch k := dest.Kind(); k {
	case reflect.Float32:
		outOfRange = floatVal > math.MaxFloat32
	case reflect.Float64:
		// No checking needed.
	}

	if outOfRange {
		return fmt.Errorf("csv: value %s is out of range for type %s", value, dest.Type().String())
	}
	dest.SetFloat(floatVal)
	return nil
}

func (u *unmarshaler) unmarshalString(dest reflect.Value, value string) error {
	dest.SetString(value)
	return nil
}

// An InvalidUnmarshalError describes an invalid argument passed to Unamrshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "csv: Unmarshal(nil)"
	}

	return "csv: Unmarshal(" + e.Type.String() + " is not a pointer to a struct pointer slice)"
}
