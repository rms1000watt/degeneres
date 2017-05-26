package generate

import (
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

	_, err := UnmarshalFile(cfg.ProtoFilePath)
	if err != nil {
		fmt.Errorf("Failed scanning protofile: %s", err)
		return
	}

	// fmt.Println(proto)
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

	for _, importFilepath := range proto.Imports {
		cnt++
		if cnt > 100 {
			fmt.Println("Greater than 100 imports.. recursive import?")
			break
		}
		filePath := filepath.Join(build.Default.GOPATH, "src", importFilepath)
		out, err := UnmarshalFile(filePath)
		if err != nil {
			fmt.Println("Failed unmarshalling file")
			break
		}

		Merge(&proto, out)
	}

	fmt.Println("Merged:", proto)
	return
}

func Merge(protoP *Proto, proto Proto) {
	// todo
}
