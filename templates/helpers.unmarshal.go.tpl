package helpers

import (
    "encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/spf13/cast"
	log "github.com/sirupsen/logrus"
)

func Unmarshal(r *http.Request, dst interface{}) (err error) {
	if r.Method == http.MethodGet {
		t := reflect.TypeOf(dst).Elem()
		v := reflect.ValueOf(dst).Elem()

		if err := r.ParseForm(); err != nil {
			return err
		}

		for i := 0; i < t.NumField(); i++ {
			jsonTag := t.Field(i).Tag.Get(TagNameJSON)
			jsonParams := strings.Split(jsonTag, ",")
			if len(jsonParams) == 0 {
				continue
			}
			jsonName := jsonParams[0]

			validateTag := t.Field(i).Tag.Get(TagNameValidate)
			validateParams := strings.Split(validateTag, ",")
			required := false
			for _, param := range validateParams {
				if param == ValidateStrRequired {
					required = true
				}
			}

			formValue := r.Form.Get(jsonName)
			if formValue == "" && required {
				return errors.New("Empty required field")
			}

			v.Field(i).Set(reflect.New(v.Field(i).Type().Elem()))

			switch v.Field(i).Type() {
			case TypeOfStringP:
				v.Field(i).Elem().SetString(formValue)
			case TypeOfIntP:
				fallthrough
			case TypeOfInt64P:
				v.Field(i).Elem().SetInt(cast.ToInt64(formValue))
			case TypeOfFloat64P:
				v.Field(i).Elem().SetFloat(cast.ToFloat64(formValue))
			case TypeOfFloat32P:
				return errors.New("Float32 not supported")
			default:
				return errors.New(fmt.Sprint("Field not set:", v.Type().Field(i).Name))
			}
		}
		return
	}

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		log.Debug("Error decoding r.Body")
		return err
	}

	return
}
