package {{.CommandLine.Command.Name}}

{{if .CommandLine.Command.API}}
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	"github.com/magical/argon2"
	"github.com/spf13/cast"
)

func Transform(in interface{}) (err error) {
	t := reflect.TypeOf(in).Elem()
	v := reflect.ValueOf(in).Elem()

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get(TagNameTransform)

		if tag == "" || tag == "-" || tag == "_" || tag == " " {
			continue
		}

		params := strings.Split(tag, ",")
		for _, param := range params {
			fmt.Printf("Transforming: %s - %s\n", v.Type().Field(i).Name, param)

			key, val := getTagKV(param)
			if v.Field(i).Pointer() == 0 && key == TransformStrDefault {
				if err := SetDefaultValue(v.Field(i), val); err != nil {
					return err
				}
				continue
			}

			if v.Field(i).Pointer() == 0 {
				return err
			}

			switch v.Field(i).Elem().Type() {
			case TypeOfString:
				if err := TransformString(param, v.Field(i).Elem()); err != nil {
					return err
				}
			}
		}
	}

	return
}

func SetDefaultValue(value reflect.Value, defaultStr string) (err error) {
	value.Set(reflect.New(value.Type().Elem()))

	switch value.Type() {
	case TypeOfStringP:
		value.Elem().SetString(defaultStr)
	case TypeOfIntP:
		value.Elem().SetInt(cast.ToInt64(defaultStr))
	case TypeOfFloat32P:
		fmt.Println("Unable to set default: Float32")
	case TypeOfFloat64P:
		value.Elem().SetFloat(cast.ToFloat64(defaultStr))
	case TypeOfBoolP:
		value.Elem().SetBool(cast.ToBool(defaultStr))
	default:
		fmt.Println("Unable to set default: no type defined")
	}
	return
}

func TransformString(param string, value reflect.Value) (err error) {
	k, v := getTagKV(param)

	switch k {
	case TransformStrHash:
		hashBytes32 := sha256.Sum256([]byte(value.String()))
		value.SetString(hex.EncodeToString(hashBytes32[:]))
	case TransformStrEncrypt:
		if value.String() == "" {
			return
		}
		if err := EncryptReflectValue(value); err != nil {
			fmt.Println("Failed Encryption...")
			return err
		}
	case TransformStrDecrypt:
		if value.String() == "" {
			return
		}
		if err := DecryptReflectValue(value); err != nil {
			fmt.Println("Failed Decryption...")
			return err
		}
	case TransformStrTrimChars:
		value.SetString(strings.Trim(value.String(), v))
	case TransformStrTrimSpace:
		value.SetString(strings.TrimSpace(value.String()))
	case TransformStrTruncate:
		truncateLength := cast.ToInt(v)
		if len(value.String()) < truncateLength {
			return
		}
		value.SetString(value.String()[:truncateLength])
	case TransformStrPasswordHash:
		if value.String() == "" {
			return
		}
		if err := PasswordHashReflectValue(value); err != nil {
			fmt.Println("Failed Password Hashing..")
			return err
		}
	}

	return
}

func EncryptReflectValue(value reflect.Value) (err error) {
	fmt.Println("DONT USE THIS KEY IN PRODUCTION.. FETCH KEY FROM PKI")
	key := []byte("AES256Key-32Characters1234567890")

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	nonce := []byte("DON'T USE ME")
	fmt.Println("DONT USE THIS NONCE IN PRODUCTION.. GENERATE AND STORE RANDOM ONE")
	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	// nonce := make([]byte, 12)
	// if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
	// 	return err
	// }

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	cipherBytes := aesgcm.Seal(nil, nonce, []byte(value.String()), nil)

	value.SetString(hex.EncodeToString(cipherBytes))
	return
}

func DecryptReflectValue(value reflect.Value) (err error) {
	fmt.Println("DONT USE THIS KEY IN PRODUCTION.. FETCH KEY FROM PKI")
	key := []byte("AES256Key-32Characters1234567890")
	ciphertext, err := hex.DecodeString(value.String())
	if err != nil {
		return err
	}

	nonce := []byte("DON'T USE ME")
	fmt.Println("DONT USE THIS NONCE IN PRODUCTION.. FETCH THE ONE FOR THIS ENTRY")

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	value.SetString(string(plaintext))
	return
}

func PasswordHashReflectValue(value reflect.Value) (err error) {
	salt, err := getRandomSalt()
	if err != nil {
		fmt.Println("Failed getting random salt:", err)
		return err
	}
	key, err := argon2.Key([]byte(value.String()), []byte(salt), 2<<14-1, 1, 8, 64)
	if err != nil {
		fmt.Println("Failed to get argon2 key:", err)
		return err
	}
	// Store these if you need to verify later
	value.SetString(hex.EncodeToString(key))
	return
}
{{end}}
