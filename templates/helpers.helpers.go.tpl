package helpers

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

const (
	TagNameValidate             = "validate"
	TagNameTransform            = "transform"
	TagNameJSON                 = "json"
	TransformStrEncrypt         = "encrypt"
	TransformStrDecrypt         = "decrypt"
	TransformStrHash            = "hash"
	TransformStrPasswordHash    = "passwordhash"
	TransformStrTruncate        = "truncate"
	TransformStrTrimChars       = "trimchars"
	TransformStrTrimSpace       = "trimspace"
	TransformStrDefault         = "default"
	ValidateStrMaxLength        = "maxlength"
	ValidateStrMinLength        = "minlength"
	ValidateStrGreaterThan      = "greaterthan"
	ValidateStrLessThan         = "lessthan"
	ValidateStrRequired         = "required"
	ValidateStrMustHaveChars    = "musthavechars"
	ValidateStrCantHaveChars    = "canthavechars"
	ValidateStrOnlyHaveChars    = "onlyhavechars"
	ValidateStrMaxLengthErr     = "Failed Max Length Validation"
	ValidateStrMinLengthErr     = "Failed Min Length Validation"
	ValidateStrRequiredErr      = "Failed Required Validation"
	ValidateStrMustHaveCharsErr = "Failed Must Have Chars Validation"
	ValidateStrCantHaveCharsErr = "Failed Can't Have Chars Validation"
	ValidateStrOnlyHaveCharsErr = "Failed Only Have Chars Validation"
	ValidateStrGreaterThanErr   = "Failed Greater Than Validation"
	ValidateStrLessThanErr      = "Failed Less Than Validation"
)

var (
	dummyString   string
	dummyInt      int
	dummyInt64    int64
	dummyFloat32  float32
	dummyFloat64  float64
	dummyBool     bool
	dummyStringP  *string
	dummyIntP     *int
	dummyInt64P   *int64
	dummyFloat32P *float32
	dummyFloat64P *float64
	dummyBoolP    *bool

	TypeOfString   = reflect.TypeOf(dummyString)
	TypeOfInt      = reflect.TypeOf(dummyInt)
	TypeOfInt64    = reflect.TypeOf(dummyInt64)
	TypeOfFloat32  = reflect.TypeOf(dummyFloat32)
	TypeOfFloat64  = reflect.TypeOf(dummyFloat64)
	TypeOfBool     = reflect.TypeOf(dummyBool)
	TypeOfStringP  = reflect.TypeOf(dummyStringP)
	TypeOfIntP     = reflect.TypeOf(dummyIntP)
	TypeOfInt64P   = reflect.TypeOf(dummyInt64P)
	TypeOfFloat32P = reflect.TypeOf(dummyFloat32P)
	TypeOfFloat64P = reflect.TypeOf(dummyFloat64P)
	TypeOfBoolP    = reflect.TypeOf(dummyBoolP)

	builtinTypes = []reflect.Type{
		TypeOfString,
		TypeOfInt,
		TypeOfInt64,
		TypeOfFloat32,
		TypeOfFloat64,
		TypeOfBool,
		TypeOfStringP,
		TypeOfIntP,
		TypeOfInt64P,
		TypeOfFloat32P,
		TypeOfFloat64P,
		TypeOfBoolP,
	}
)

func getRandomSalt() (salt []byte, err error) {
	salt = make([]byte, 32)
	_, err = rand.Read(salt)
	return
}

func getTagKV(param string) (k, v string) {
	paramArr := strings.Split(param, "=")

	k = paramArr[0]
	if len(paramArr) == 2 {
		v = paramArr[1]
	}
	k = strings.ToLower(k)
	k = strings.Replace(k, "-", "", -1)
	k = strings.Replace(k, "_", "", -1)
	k = strings.Replace(k, " ", "", -1)
	return
}

func allCharsInStr(allChars, in string) (out bool) {
	for _, char := range allChars {
		if strings.Index(in, string(char)) == -1 {
			return
		}
	}
	return true
}

func onlyCharsInStr(onlyChars, in string) (out bool) {
	for _, char := range onlyChars {
		in = strings.Replace(in, string(char), "", -1)
	}
	return len(in) == 0
}

func dereferenceStringArray(in []*string) (out []string) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

func dereferenceIntArray(in []*int) (out []int) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

func dereferenceInt32Array(in []*int32) (out []int32) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

func dereferenceInt64Array(in []*int64) (out []int64) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

func dereferenceFloat32Array(in []*float32) (out []float32) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

func dereferenceFloat64Array(in []*float64) (out []float64) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

func dereferenceBoolArray(in []*bool) (out []bool) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

func isBuiltin(fieldType reflect.Type) bool {
	for _, builtinType := range builtinTypes {
		if fieldType == builtinType {
			return true
		}
	}
	
	if strings.Contains(fieldType.String(), "map[") {
		return true
	}

	ft := strings.Replace(fieldType.String(), "[]*", "", -1)
	for _, builtinType := range builtinTypes {
		if ft == builtinType.String() {
			return true
		}
	}

	return false
}

// Courtesy of https://stackoverflow.com/questions/13901819/quick-way-to-detect-empty-values-via-reflection-in-go
func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func ExecDirectory(commandStr string, directory string, envVars ...string) (err error) {
	if err := os.Chdir(directory); err != nil {
		log.Errorf("Failed to cd %s: %s", directory, err)
		return err
	}

	return Exec(commandStr, envVars...)
}

func Exec(commandStr string, envVars ...string) (err error) {
	if strings.TrimSpace(commandStr) == "" {
		return errors.New("No command provided")
	}

	var name string
	var args []string

	cmdArr := strings.Split(commandStr, " ")
	name = cmdArr[0]

	if len(cmdArr) > 1 {
		args = cmdArr[1:]
	}

	command := exec.Command(name, args...)
	command.Env = append(os.Environ(), envVars...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		log.Error("Failed creating command stdoutpipe: ", err)
		return err
	}
	defer stdout.Close()
	stdoutReader := bufio.NewReader(stdout)

	stderr, err := command.StderrPipe()
	if err != nil {
		log.Error("Failed creating command stderrpipe: ", err)
		return err
	}
	defer stderr.Close()
	stderrReader := bufio.NewReader(stderr)

	if err := command.Start(); err != nil {
		log.Error("Failed starting command: ", err)
		return err
	}

	go handleReader(stdoutReader, false)
	go handleReader(stderrReader, true)

	if err := command.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Debug("Exit Status: ", status.ExitStatus())
				return err
			}
		}
		log.Debug("Failed to wait for command: ", err)
		return err
	}

	return
}

func handleReader(reader *bufio.Reader, isStderr bool) {
	printOutput := log.GetLevel() == log.DebugLevel
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if printOutput {
			fmt.Print(str)
		}
	}
}

