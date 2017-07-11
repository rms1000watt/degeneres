//go:generate rm -rf data helpers
//go:generate go run $(go env GOPATH)/src/github.com/rms1000watt/degeneres/main.go generate -f {{.ProtoFilePath}} -o $(pwd)

package main
