package generate

import (
	"encoding/json"
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"
)

var (
	cnt = 0
)

func Generate(cfg Config) {
	fmt.Println("Starting Generation...")

	proto, err := UnmarshalFile(cfg.ProtoFilePath)
	if err != nil {
		fmt.Errorf("Failed scanning protofile: %s", err)
		return
	}

	NewTemplate(proto)

	bytes, _ := json.MarshalIndent(proto, "", "  ")
	fmt.Println(string(bytes))
}

func UnmarshalFile(filePath string) (proto Proto, err error) {
	fmt.Println("Starting Unmarshal...", filePath)

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Errorf("Failed reading file: %s: %s", filePath, err)
		return
	}

	tokens := Scan(fileBytes)
	proto = Parse(tokens)

	importedProtos := []Proto{}
	for _, importFilepath := range proto.Imports {
		cnt++
		if cnt > 100 {
			fmt.Println("Greater than 100 imports.. recursive import?")
			break
		}
		filePath := filepath.Join(build.Default.GOPATH, "src", importFilepath)
		importedProto, err := UnmarshalFile(filePath)
		if err != nil {
			fmt.Println("Failed unmarshalling file:", err)
			break
		}

		importedProtos = append(importedProtos, importedProto)
	}

	if err := Merge(&proto, importedProtos...); err != nil {
		return proto, err
	}

	return
}
