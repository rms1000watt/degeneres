package generate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	TagNameValidate  = "validate"
	ValidateRequired = "required"
)

var (
	ValidateRequiredErr = errors.New("Failed Required Validation")
)

func Validate(in interface{}) (err error) {
	t := reflect.TypeOf(in).Elem()
	v := reflect.ValueOf(in).Elem()

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get(TagNameValidate)

		if tag == "" || tag == "-" || tag == "_" || tag == " " {
			continue
		}

		params := strings.Split(tag, ",")
		for _, param := range params {
			fmt.Printf("Validating: %s - %s\n", v.Type().Field(i).Name, param)

			if param == ValidateRequired {
				if v.Field(i).String() == "" {
					return ValidateRequiredErr
				}
			}
		}
	}
	return
}
