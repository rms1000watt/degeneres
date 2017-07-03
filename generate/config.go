package generate

import (
	log "github.com/sirupsen/logrus"
)

type Config struct {
	ProtoFilePath string
	Verbose       bool
	OutPath       string
	LogLevel      log.Level
}
