package generate

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Certs(certsPath, commonName string) {
	fmt.Println("Generating certs...")

	if err := os.Chdir(certsPath); err != nil {
		fmt.Println("Failed chdir to certsPath:", certsPath)
		return
	}

	keys := []string{
		"ca.cer",
		"ca.key",
		"server.csr",
		"server.key",
		"server.cer",
	}
	for _, key := range keys {
		os.Remove(key)
	}

	// # Courtesy of https://github.com/deckarep/EasyCert
	cmds := []string{
		"openssl genrsa -out ca.key 2048",
		`openssl req -x509 -new -key ca.key -out ca.cer -days 90 -subj /CN="rms1000watt-CA"`,
		"openssl genrsa -out server.key 2048",
		`openssl req -new -key server.key -out server.csr -config ./openssl.cnf -subj /CN="` + commonName + `"`,
		"openssl x509 -req -in server.csr -out server.cer -CAkey ca.key -CA ca.cer -days 90 -CAcreateserial -CAserial serial",
	}

	for _, cmd := range cmds {
		if err := execute(cmd); err != nil {
			return
		}
	}
}

func execute(cmd string) (err error) {
	cmdArr := strings.Split(cmd, " ")
	if len(cmdArr) < 2 {
		return errors.New("bad command provided")
	}

	name := cmdArr[0]
	args := cmdArr[1:]

	outBytes, err := exec.Command(name, args...).CombinedOutput()
	fmt.Println(string(outBytes))

	if err != nil {
		fmt.Printf("Failed executing cmd: '%s': %s\n", cmd, err)
		return err
	}

	return
}
