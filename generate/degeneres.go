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
	Services             []DgService // Commands.. go run main.go $SERVICE_NAME
	Messages             []DgMessage
}

type DgService struct {
	Name
	ShortDescription string
	LongDescription  string
	Middlewares      map[string]DgMiddleware
	Endpoints        []DgEndpoint
	Messages         []DgMessage

	CertsPath      string
	ImportPath     string
	PublicKeyName  string
	PrivateKeyName string
}

type DgMiddleware struct {
	Options []KV
}

type DgEndpoint struct {
	Name
	Pattern     string
	Middlewares map[string]DgMiddleware
	Methods     []string
	In          string
	Out         string
}

type DgMessage struct {
	Name
	Fields  []DgField
	IsInput bool
}

type DgField struct {
	Name
	DataType  string // Process all of these to include pointers if it's an input
	Transform string // Parse from Options
	Validate  string // Parse from Options
}

type Name struct {
	Raw        string
	Dash       string
	Snake      string
	Camel      string
	Lower      string
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
				Name:      genName(protoField.Name),
				DataType:  fixDataType(protoField.DataType, true, protoField.Rule),
				Transform: getTransformFromOptions(protoField.Options),
				Validate:  getValidateFromOptions(protoField.Options),
			})

		}

		messages = append(messages, DgMessage{
			Name:   genName(protoMessage.Name),
			Fields: fields,
		})
	}

	dg.Messages = messages

	services := []DgService{}
	for _, service := range proto.Services {
		spew.Dump(genName(service.Name))

		// endpoints := []DgEndpoint{}

		services = append(services, DgService{
			Name:           genName(service.Name),
			CertsPath:      dg.CertsPath,
			ImportPath:     dg.ImportPath,
			PublicKeyName:  dg.PublicKeyName,
			PrivateKeyName: dg.PrivateKeyName,
			Messages:       messages,
		})
	}
	dg.Services = services

	err = Validate(&dg)

	return
}

func fixOptionName(in string) (out string) {
	inArr := strings.Split(in, ".")
	if len(inArr) < 2 {
		return in
	}

	return inArr[1]
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

func fixDataType(dataType string, isInput bool, fieldRule string) string {
	isRepeated := strings.ToLower(fieldRule) == FieldRuleRepeated

	splitDT := strings.Split(dataType, ".")
	if len(splitDT) > 1 {
		dataType = splitDT[len(splitDT)-1]
	}

	if isInput {
		if len(dataType) > 2 && isRepeated {
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
