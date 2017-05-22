package generate

import (
	"fmt"
)

func Generate(cfg Config) {
	fmt.Println("Starting Generation...")

	_, err := Scan(cfg.ProtoFilePath)
	if err != nil {
		fmt.Errorf("Failed scanning protofile: %s", err)
		return
	}
}
