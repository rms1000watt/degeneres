package helpers

import (
	"reflect"
	"strings"

	"github.com/spf13/cast"
	log "github.com/sirupsen/logrus"
)

// TODO: Change with https://play.golang.org/p/u4UQbo0ZeM
func Validate(in interface{}) (msg string, err error) {
	t := reflect.TypeOf(in).Elem()
	v := reflect.ValueOf(in).Elem()

	for i := 0; i < t.NumField(); i++ {
		if !isBuiltin(t.Field(i).Type) {
			if IsZeroOfUnderlyingType(v.Field(i).Interface()) {
				continue
			}

			msg, err := Validate(v.Field(i).Interface())
			if err != nil {
				log.Debug("Error field validate:", err)
				return msg, err
			}

			if msg != "" {
				log.Debug("Failed field validate:", msg)
				return msg, err
			}
		}

		tag := t.Field(i).Tag.Get(TagNameValidate)

		if tag == "" || tag == "-" || tag == "_" || tag == " " {
			continue
		}

		fieldPointer := v.Field(i).Pointer()
		if strings.Contains(strings.ToLower(tag), ValidateStrRequired) {
			if fieldPointer == 0 {
				return ValidateStrRequiredErr, nil
			}
		}

		if fieldPointer == 0 {
			continue
		}

		params := strings.Split(tag, ",")
		for _, param := range params {
			log.Debugf("Validating: %s - %s", v.Type().Field(i).Name, param)

			switch v.Field(i).Elem().Type() {
			case TypeOfString:
				if vMsg := ValidateString(param, v.Field(i).Elem().String()); vMsg != "" {
					return vMsg, nil
				}
			case TypeOfInt:
				if vMsg := ValidateInt(param, int(v.Field(i).Elem().Int())); vMsg != "" {
					return vMsg, nil
				}
			case TypeOfFloat64:
				if vMsg := ValidateFloat64(param, v.Field(i).Elem().Float()); vMsg != "" {
					return vMsg, nil
				}
			}
		}
	}

	return
}

func ValidateString(param, in string) (msg string) {
	k, v := getTagKV(param)

	switch k {
	case ValidateStrMaxLength:
		if len(in) > cast.ToInt(v) {
			return ValidateStrMaxLengthErr
		}
	case ValidateStrMinLength:
		if len(in) < cast.ToInt(v) {
			return ValidateStrMinLengthErr
		}
	case ValidateStrMustHaveChars:
		if !allCharsInStr(v, in) {
			return ValidateStrMustHaveCharsErr
		}
	case ValidateStrCantHaveChars:
		if strings.IndexAny(in, v) > -1 {
			return ValidateStrCantHaveCharsErr
		}
	case ValidateStrOnlyHaveChars:
		if !onlyCharsInStr(v, in) {
			return ValidateStrOnlyHaveCharsErr
		}
	}

	return
}

func ValidateInt(param string, in int) (msg string) {
	k, v := getTagKV(param)

	switch k {
	case ValidateStrGreaterThan:
		if in < cast.ToInt(v) {
			return ValidateStrGreaterThanErr
		}
	case ValidateStrLessThan:
		if in > cast.ToInt(v) {
			return ValidateStrLessThanErr
		}
	}

	return
}

func ValidateFloat64(param string, in float64) (msg string) {
	k, v := getTagKV(param)

	switch k {
	case ValidateStrGreaterThan:
		if in < cast.ToFloat64(v) {
			return ValidateStrGreaterThanErr
		}
	case ValidateStrLessThan:
		if in > cast.ToFloat64(v) {
			return ValidateStrLessThanErr
		}
	}

	return
}
