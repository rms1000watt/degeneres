package generate

import (
	"errors"
	"strings"
)

var (
	ErrorMismatchSyntax     = errors.New("Syntax mismatch on imports")
	ErrorNoImportsProvided  = errors.New("No imports provided")
	ErrorFieldInvalidSyntax = errors.New("Field import syntax invalid")
)

func Merge(proto *Proto, importedProtos ...Proto) (err error) {
	importFields := []string{}
	for _, message := range proto.Messages {
		for _, field := range message.Fields {
			if strings.Contains(field.DataType, ".") {
				importFields = append(importFields, field.DataType)
			}
		}
	}

	if len(importFields) == 0 {
		return
	}

	if len(importFields) > 0 && len(importedProtos) == 0 {
		return ErrorNoImportsProvided
	}

	importedProtosMap := map[string]Proto{}
	for _, importedProto := range importedProtos {
		importedProtosMap[importedProto.Package] = importedProto
	}

	for _, importField := range importFields {
		fieldSplit := strings.Split(importField, ".")
		if len(fieldSplit) == 0 {
			return ErrorFieldInvalidSyntax
		}
		importPackage := fieldSplit[0]
		importProto := importedProtosMap[importPackage]
		for _, message := range importProto.Messages {
			if message.Name == fieldSplit[1] {
				message.Imported = true
				proto.Messages = append(proto.Messages, message)
			}
		}
	}

	for _, importedProto := range importedProtos {
		for _, message := range importedProto.Messages {
			if message.Imported {
				proto.Messages = append(proto.Messages, message)
			}
		}
	}

	return
}
