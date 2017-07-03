package generate

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	dirTemplates = "templates"
	dirHelpers   = "helpers"
	extTpl       = ".tpl"
)

var (
	cnt        = 0
	errGenFail = errors.New("Failed generating project")
	funcMap    = template.FuncMap{
		"TimeNowYear":          time.Now().Year,
		"HandleQuotes":         HandleQuotes,
		"EmptyValue":           EmptyValue,
		"GetHTTPMethod":        GetHTTPMethod,
		"FallbackSet":          FallbackSet,
		"GetMethodMiddlewares": GetMethodMiddlewares,
		"GetPathMiddlewares":   GetPathMiddlewares,
		"GetInputType":         GetInputType,
		"GetDereferenceFunc":   GetDereferenceFunc,
		"IsStruct":             IsStruct,
		"MinusP":               MinusP,
		"MinusStar":            MinusStar,
	}
)

func Generate(cfg Config) {
	log.Debug("Starting generator")
	defer log.Debug("Generator done")

	proto, err := UnmarshalFile(cfg.ProtoFilePath)
	if err != nil {
		log.Error("Failed scanning protofile: ", err)
		return
	}

	dg, err := NewDegeneres(proto)
	if err != nil {
		log.Error("Failed converting to degeneres format: ", err)
		return
	}

	dg.GeneratorVersion = getGeneratorVersion()

	if err := os.Mkdir(cfg.OutPath, os.ModePerm); err != nil {
		log.Errorf("Directory: \"%s\" already exists. Continuing...\n", cfg.OutPath)
	}

	helperFileNames, err := getHelperFileNames()
	if err != nil {
		log.Error("Failed reading helper files: ", err)
		return
	}

	templates := getTemplates(dg)
	for _, tpl := range templates {
		if err := genFile(cfg, tpl, helperFileNames); err != nil {
			log.Error("Failed generating template: ", tpl.TemplateName, ": ", err)
		}
	}
}

func UnmarshalFile(filePath string) (proto Proto, err error) {
	log.Debug("Starting unmarshal: ", filePath)
	defer log.Debug("Unmarshal done: ", filePath)

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("Failed reading file: %s: %s", filePath, err)
		return
	}

	tokens := Scan(fileBytes)
	proto = Parse(tokens)

	importedProtos := []Proto{}
	for _, importFilepath := range proto.Imports {
		cnt++
		if cnt > 100 {
			log.Warn("Greater than 100 imports.. recursive import?")
			break
		}
		filePath := filepath.Join(build.Default.GOPATH, "src", importFilepath)
		importedProto, err := UnmarshalFile(filePath)
		if err != nil {
			log.Error("Failed unmarshalling file: ", err)
			break
		}

		importedProtos = append(importedProtos, importedProto)
	}

	if err := Merge(&proto, importedProtos...); err != nil {
		return proto, err
	}

	return
}

func getHelperFileNames() (helperFileNames []string, err error) {
	helperFiles, err := ioutil.ReadDir(filepath.Join(build.Default.GOPATH, "src", "github.com", "rms1000watt", "degeneres", dirTemplates, dirHelpers))
	if err != nil {
		return
	}

	for _, helperFile := range helperFiles {
		helperFileNames = append(helperFileNames, filepath.Join(dirHelpers, helperFile.Name()))
	}
	return
}

func getTemplates(dg Degeneres) (templates []Template) {
	singleTemplateNames := []string{
		".gitignore.tpl",
		"main.go.tpl",
		"cmd.root.go.tpl",
		"cmd.version.go.tpl",
		"Readme.md.tpl",
		"License..tpl",
		"Dockerfile..tpl",
		"helpers.middlewares.go.tpl",
		"helpers.helpers.go.tpl",
		"helpers.unmarshal.go.tpl",
		"helpers.validate.go.tpl",
		"helpers.transform.go.tpl",
		"helpers.handler.go.tpl",
		"data.data.go.tpl",
		"data.input.go.tpl",
	}

	for _, singleTemplateName := range singleTemplateNames {
		templates = append(templates, Template{
			TemplateName: singleTemplateName,
			FileName:     singleTemplateName,
			Data:         dg,
		})
	}

	for _, service := range dg.Services {
		lowerKey := service.Camel

		templates = append(templates, Template{
			TemplateName: "cmd." + lowerKey + ".go.tpl",
			FileName:     "cmd.command.go.tpl",
			Data:         service,
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + "." + lowerKey + ".go.tpl",
			FileName:     "command.command.go.tpl",
			Data:         service,
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + ".config.go.tpl",
			FileName:     "command.config.go.tpl",
			Data:         service,
		})

		for _, endpoint := range service.Endpoints {
			// ephemeralCfg := templateCfg
			// ephemeralCfg.API.Paths = []TemplatePath{apiPath}

			templates = append(templates, Template{
				TemplateName: fmt.Sprintf("%s.%sHandler.go.tpl", lowerKey, endpoint.Camel),
				FileName:     "command.handler.go.tpl",
				Data:         endpoint,
				// Data: service,
			})
		}
	}

	return
}

func genFile(cfg Config, tpl Template, helperFileNames []string) (err error) {
	templateName := tpl.TemplateName

	templateNameArr := strings.Split(templateName, ".")
	if len(templateNameArr) < 3 {
		log.Error("Bad templateName provided: ", templateName)
		return errGenFail
	}

	templateFileName := filepath.Join(build.Default.GOPATH, "src", "github.com", "rms1000watt", "degeneres", dirTemplates, tpl.FileName)
	fileBytes, err := ioutil.ReadFile(templateFileName)
	if err != nil {
		log.Error("Failed reading template file: ", err)
		return
	}

	t, err := template.New(templateName).Funcs(funcMap).Parse(string(fileBytes))
	if err != nil {
		log.Error("Failed parsing template: ", err)
		return errGenFail
	}

	fullHelperFileNames := []string{}
	for _, helperFileName := range helperFileNames {
		fullHelperFileNames = append(fullHelperFileNames, filepath.Join(build.Default.GOPATH, "src", "github.com", "rms1000watt", "degeneres", dirTemplates, helperFileName))
	}

	if _, err := t.ParseFiles(fullHelperFileNames...); err != nil {
		log.Error("Failed parsing template: ", err)
		return errGenFail
	}

	var outBuffer bytes.Buffer
	if err := t.Execute(&outBuffer, tpl.Data); err != nil {
		log.Error("Failed executing template: ", err)
		return errGenFail
	}

	// Make the required directories for the project
	dirs := templateNameArr[:len(templateNameArr)-3]
	dirs = append([]string{cfg.OutPath}, dirs...)

	if len(dirs) != 0 {
		if err := os.MkdirAll(filepath.Join(dirs...), os.ModePerm); err != nil {
			log.Error("Failed mkdir on dirs: ", dirs, ": ", err)
			return errGenFail
		}
	}

	dirsStr := filepath.Join(dirs...)
	fileName := strings.Join(templateNameArr[len(templateNameArr)-3:len(templateNameArr)-1], ".")
	completeFilePath := filepath.Join(dirsStr, fileName)

	if filepath.Ext(completeFilePath) == "." {
		completeFilePath = completeFilePath[:len(completeFilePath)-1]
	}

	if _, err := os.Stat(completeFilePath); err == nil {
		log.Info("NO overwrite, file exists: ", completeFilePath)
		return nil
	}

	log.Info("Writing: ", completeFilePath)
	if err := ioutil.WriteFile(completeFilePath, outBuffer.Bytes(), os.ModePerm); err != nil {
		log.Error("Failed writing file: ", err)
		return errGenFail
	}

	if filepath.Ext(completeFilePath) == ".go" {
		exec.Command("goimports", "-w", completeFilePath).CombinedOutput()
		exec.Command("gofmt", "-w", completeFilePath).CombinedOutput()
	}

	// TODO: Remove this--be smarter about which files to write
	RemoveUnusedFile(completeFilePath)

	return nil
}

func RemoveUnusedFile(completeFilePath string) {
	fileBytes, err := ioutil.ReadFile(completeFilePath)
	if err != nil {
		// Fail silently.. not a big deal
		return
	}

	// if !bytes.Contains(bytes.TrimSpace(fileBytes), []byte("\n")) && bytes.Equal(fileBytes[:7], []byte("package")) {
	if !bytes.Contains(bytes.TrimSpace(fileBytes), []byte("\n")) {
		fmt.Println("Removing:", completeFilePath)
		if err := os.Remove(completeFilePath); err != nil {
			// Fail silently.. not a big deal
			return
		}
	}
}

func getGeneratorVersion() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		return "Not sure: Broken git"
	}
	return strings.TrimSpace(string(out))
}
