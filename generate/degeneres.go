package generate

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/davecgh/go-spew/spew"
)

const (
	OptionVersion              = "version"
	OptionImportPath           = "import_path"
	OptionAuthor               = "author"
	OptionProjectName          = "project_name"
	OptionShortDescription     = "short_description"
	OptionLongDescription      = "long_description"
	OptionCertsPath            = "certs_path"
	OptionPublicKeyName        = "public_key_name"
	OptionPrivateKeyName       = "private_key_name"
	OptionDockerPath           = "docker_path"
	OptionProjectNameCommander = "project_name_commander"
	OptionTransform            = "transform"
	OptionValidate             = "validate"
	OptionMiddlewareCORS       = "middleware.cors"
	OptionMiddlewareNoCache    = "middleware.no_cache"
	OptionMethod               = "method"

	MiddlewareCORS    = "CORS"
	MiddlewareNoCache = "NoCache"
)

var (
	protoDataTypes = []string{
		"double",
		"float",
		"int32",
		"int64",
		"uint32",
		"uint64",
		"sint32",
		"sint64",
		"fixed32",
		"fixed64",
		"sfixed32",
		"sfixed64",
		"bool",
		"string",
		"byte",
	}
)

type Degeneres struct {
	Version              string `validate:"required"`
	ImportPath           string `validate:"required"`
	DockerPath           string
	Author               string `validate:"required"`
	ProjectName          string `validate:"required"`
	ProjectNameCommander string
	ProjectFolder        string
	ShortDescription     string
	LongDescription      string
	CertsPath            string
	PublicKeyName        string
	PrivateKeyName       string
	Services             []DgService
	Messages             []DgMessage
	Inputs               []DgMessage
}

type DgService struct {
	Name
	ShortDescription string
	LongDescription  string
	Middlewares      map[string]string
	MiddlewareNames  string
	Endpoints        []DgEndpoint

	CertsPath      string
	ImportPath     string
	PublicKeyName  string
	PrivateKeyName string
}

type DgEndpoint struct {
	Name
	ServiceName     Name
	Pattern         string
	Middlewares     map[string]string
	MiddlewareNames string
	Methods         []Name
	Input           Name
	Output          Name

	ImportPath string
}

type DgMessage struct {
	Name
	Fields  []DgField
	IsInput bool
}

type DgField struct {
	Name
	DataType         string
	Transform        string
	Validate         string
	Rule             string
	IsRepeated       bool
	IsStruct         bool
	IsRepeatedStruct bool
}

type Name struct {
	Raw        string
	Dash       string
	Snake      string
	Camel      string
	Lower      string
	LowerDash  string
	LowerSnake string
	LowerCamel string
	Title      string
	TitleSnake string
	TitleCamel string
	Upper      string
	UpperSnake string
	UpperCamel string
}

func NewDegeneres(proto Proto) (dg Degeneres, err error) {
	dg = Degeneres{}
	for _, option := range proto.Options {
		optionName := strings.ToLower(fixOptionName(option.Name))
		switch optionName {
		case OptionAuthor:
			dg.Author = option.Value
		case OptionCertsPath:
			dg.CertsPath = option.Value
		case OptionImportPath:
			dg.ImportPath = option.Value
			splitPath := strings.Split(dg.ImportPath, "/")
			dg.ProjectFolder = splitPath[len(splitPath)-1]
		case OptionLongDescription:
			dg.LongDescription = option.Value
		case OptionShortDescription:
			dg.ShortDescription = option.Value
		case OptionPrivateKeyName:
			if option.Value == "" {
				dg.PrivateKeyName = "server.key"
				continue
			}
			dg.PrivateKeyName = option.Value
		case OptionPublicKeyName:
			if option.Value == "" {
				dg.PublicKeyName = "server.cer"
				continue
			}
			dg.PublicKeyName = option.Value
		case OptionProjectName:
			dg.ProjectName = option.Value
			dg.ProjectNameCommander = ToDashCase(option.Value)
		case OptionVersion:
			dg.Version = option.Value
		case OptionDockerPath:
			dg.DockerPath = option.Value
		}
	}

	messages := []DgMessage{}
	for _, protoMessage := range proto.Messages {
		fields := []DgField{}
		for _, protoField := range protoMessage.Fields {
			fields = append(fields, DgField{
				Name:             genName(protoField.Name),
				DataType:         fixDataType(protoField.DataType, false, protoField.Rule),
				Transform:        getTransformFromOptions(protoField.Options),
				Validate:         getValidateFromOptions(protoField.Options),
				Rule:             protoField.Rule,
				IsRepeated:       getIsRepeated(protoField.Rule, protoField.DataType),
				IsStruct:         getIsStruct(protoField.DataType),
				IsRepeatedStruct: getIsRepeatedStruct(protoField.Rule, protoField.DataType),
			})

		}

		message := DgMessage{
			Name:   genName(protoMessage.Name),
			Fields: fields,
		}

		if !messageInMessages(message, messages) {
			messages = append(messages, message)
		}
	}
	dg.Messages = messages

	inputs := []DgMessage{}
	additionalInputs := []string{}
	for ind, message := range messages {
		for _, service := range proto.Services {
			for _, rpc := range service.RPCs {
				if message.Raw == rpc.Input {
					messages[ind].IsInput = true

					fields := []DgField{}
					for _, field := range message.Fields {

						// TODO: Get additional inputs
						isStruct := getIsStruct(field.DataType)
						if isStruct {
							additionalInputs = append(additionalInputs, field.Raw)
						}

						fields = append(fields, DgField{
							Name:             genName(field.Raw),
							DataType:         fixDataType(field.DataType, true, field.Rule),
							Transform:        field.Transform,
							Validate:         field.Validate,
							Rule:             field.Rule,
							IsRepeated:       getIsRepeated(field.Rule, field.DataType),
							IsStruct:         isStruct,
							IsRepeatedStruct: getIsRepeatedStruct(field.Rule, field.DataType),
						})
					}
					input := DgMessage{
						Name:   genName(message.Raw + "P"),
						Fields: fields,
					}
					if !messageInMessages(input, inputs) {
						inputs = append(inputs, input)
					}
				}
			}
		}
	}
	// TODO: Handle additional inputs
	for _, message := range messages {
		for _, additionalInput := range additionalInputs {
			if additionalInput == message.Raw {
				fields := []DgField{}
				for _, field := range message.Fields {

					isStruct := getIsStruct(field.DataType)
					if isStruct {
						additionalInputs = append(additionalInputs, field.Raw)
					}

					fields = append(fields, DgField{
						Name:             genName(field.Raw),
						DataType:         fixDataType(field.DataType, true, field.Rule),
						Transform:        field.Transform,
						Validate:         field.Validate,
						Rule:             field.Rule,
						IsRepeated:       getIsRepeated(field.Rule, field.DataType),
						IsStruct:         isStruct,
						IsRepeatedStruct: getIsRepeatedStruct(field.Rule, field.DataType),
					})
				}
				input := DgMessage{
					Name:   genName(message.Raw + "P"),
					Fields: fields,
				}
				if !messageInMessages(input, inputs) {
					inputs = append(inputs, input)
				}
			}
		}
	}

	dg.Inputs = inputs

	services := []DgService{}
	for _, service := range proto.Services {
		spew.Dump(genName(service.Name))

		longDescription := ""
		shortDescription := ""
		for _, option := range service.Options {
			optionName := strings.ToLower(fixOptionName(option.Name))
			switch optionName {
			case OptionShortDescription:
				shortDescription = option.Value
			case OptionLongDescription:
				longDescription = option.Value
			}
		}

		middlewares, middlewareNames := getMiddlewares(service.Options)

		endpoints := []DgEndpoint{}
		for _, rpc := range service.RPCs {
			rpcMws, rpcMwNames := getMiddlewares(rpc.Options)

			endpointName := genName(rpc.Name)
			endpoints = append(endpoints, DgEndpoint{
				Name:            endpointName,
				ServiceName:     genName(service.Name),
				Pattern:         "/" + endpointName.LowerDash,
				Middlewares:     rpcMws,
				MiddlewareNames: rpcMwNames,
				Methods:         getMethods(rpc.Options),
				Input:           genName(rpc.Input),
				Output:          genName(rpc.Output),

				ImportPath: dg.ImportPath,
			})
		}

		services = append(services, DgService{
			Name:             genName(service.Name),
			ShortDescription: shortDescription,
			LongDescription:  longDescription,
			Middlewares:      middlewares,
			MiddlewareNames:  middlewareNames,
			Endpoints:        endpoints,

			CertsPath:      dg.CertsPath,
			ImportPath:     dg.ImportPath,
			PublicKeyName:  dg.PublicKeyName,
			PrivateKeyName: dg.PrivateKeyName,
		})
	}
	dg.Services = services

	err = Validate(&dg)

	return
}

func getMiddlewares(options []Option) (map[string]string, string) {
	middlewareNameArr := []string{}
	middlewares := map[string]string{}
	for _, option := range options {
		optionName := strings.ToLower(fixOptionName(option.Name))
		middlewareName := Name{}
		switch optionName {
		case OptionMiddlewareCORS:
			middlewareName = genName(MiddlewareCORS)
		case OptionMiddlewareNoCache:
			middlewareName = genName(MiddlewareNoCache)
		default:
			continue
		}
		middlewares[middlewareName.TitleCamel] = option.Value
		middlewareNameArr = append(middlewareNameArr, "helpers.Middleware"+middlewareName.TitleCamel)
	}

	return middlewares, strings.Join(middlewareNameArr, ", ")
}

func getMethods(options []Option) []Name {
	methods := []Name{}
	for _, option := range options {
		optionName := strings.ToLower(fixOptionName(option.Name))
		switch optionName {
		case OptionMethod:
			methods = append(methods, genName(strings.ToLower(option.Value)))
		}
	}
	return methods
}

func fixOptionName(in string) (out string) {
	in = strings.TrimSpace(in)
	inArr := strings.Split(in, ".")
	if len(inArr) < 2 {
		return in
	}

	return strings.Join(inArr[1:], ".")
}

func genName(in string) Name {
	camel := ToCamelCase(in)
	snake := ToSnakeCase(in)
	dash := ToDashCase(in)

	return Name{
		Raw:        in,
		Dash:       dash,
		Camel:      camel,
		Snake:      snake,
		Lower:      strings.ToLower(in),
		LowerDash:  strings.ToLower(dash),
		LowerSnake: strings.ToLower(snake),
		LowerCamel: strings.ToLower(camel),
		Upper:      strings.ToUpper(in),
		UpperSnake: strings.ToUpper(snake),
		UpperCamel: strings.ToUpper(camel),
		Title:      strings.Title(in),
		TitleSnake: strings.Title(snake),
		TitleCamel: strings.Title(camel),
	}
}

func getTransformFromOptions(options []Option) string {
	for _, option := range options {
		optionName := strings.ToLower(fixOptionName(option.Name))
		if optionName == OptionTransform {
			return option.Value
		}
	}
	return ""
}

func getValidateFromOptions(options []Option) string {
	for _, option := range options {
		optionName := strings.ToLower(fixOptionName(option.Name))
		if optionName == OptionValidate {
			return option.Value
		}
	}
	return ""
}

func getIsRepeated(fieldRule, dataType string) bool {
	return strings.ToLower(fieldRule) == FieldRuleRepeated && !getIsStruct(dataType)
}

func getIsStruct(dataType string) bool {
	dataType = strings.ToLower(dataType)
	for _, dt := range protoDataTypes {
		if dataType == dt {
			return false
		}
	}
	return true
}

func getIsRepeatedStruct(fieldRule, dataType string) bool {
	return strings.ToLower(fieldRule) == FieldRuleRepeated && getIsStruct(dataType)
}

func messageInMessages(message DgMessage, messages []DgMessage) bool {
	for _, knownMessage := range messages {
		if knownMessage.Raw == message.Raw {
			return true
		}
	}
	return false
}

func fixDataType(dataType string, isInput bool, fieldRule string) string {
	isRepeated := strings.ToLower(fieldRule) == FieldRuleRepeated

	splitDT := strings.Split(dataType, ".")
	if len(splitDT) > 1 {
		dataType = splitDT[len(splitDT)-1]
	}

	if isRepeated && len(dataType) > 2 && dataType[:2] != "[]" {
		dataType = "[]" + dataType
	}

	if isInput {
		if len(dataType) > 2 && isRepeated && dataType[:2] == "[]" {
			return "[]*" + dataType[2:]
		}
		if len(dataType) > 2 && isRepeated && dataType[:2] != "[]" {
			return "[]*" + dataType
		}
		return "*" + dataType
	}
	return dataType
}

// Courtesy of https://github.com/etgryphon/stringUp/blob/master/stringUp.go
var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

func ToCamelCase(src string) (out string) {
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	out = string(bytes.Join(chunks, nil))
	out = strings.ToLower(string(out[0])) + string(out[1:])
	return out
}

// Courtesy of https://github.com/fatih/camelcase/blob/master/camelcase.go
func ToSnakeCase(src string) (out string) {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return src
	}
	entries := []string{}
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 && !strings.Contains(string(s), " ") {
			entries = append(entries, string(s))
		}
	}

	out = strings.ToLower(strings.Join(entries, "_"))

	for strings.Contains(out, "__") {
		out = strings.Replace(out, "__", "_", -1)
	}

	return out
}

func ToDashCase(in string) (out string) {
	return strings.Replace(ToSnakeCase(in), "_", "-", -1)
}
