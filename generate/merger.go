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
		for fieldInd, field := range message.Fields {
			if strings.Contains(field.DataType, ".") {
				importFields = append(importFields, field.DataType)
				continue
			}

			// TODO: Consider fixing this hack?
			// Doesn't contain "." and is not builtIn.. Assume its a self reference?
			if !builtIn(field.DataType) {
				options := []Option{}
				for _, option := range field.Options {
					options = append(options, Option{
						Name:  option.Name,
						Value: option.Value,
					})
				}

				dataType := proto.Package + "." + field.DataType

				message.Fields[fieldInd] = Field{
					Name:     field.Name,
					DataType: dataType,
					Position: field.Position,
					Rule:     field.Rule,
					Options:  options,
				}

				importFields = append(importFields, dataType)
			}
		}
	}

	if len(importFields) == 0 {
		return
	}

	// TODO: Consider fixing this hack?
	for _, importField := range importFields {
		if strings.Contains(importField, proto.Package) && len(importedProtos) == 0 {
			importedProtos = append(importedProtos, *proto)
		}
	}

	if len(importFields) > 0 && len(importedProtos) == 0 {
		return ErrorNoImportsProvided
	}

	importedProtosMap := map[string]Proto{}
	for _, importedProto := range importedProtos {
		importedProtosMap[importedProto.Package] = importedProto
		proto.ProtoPaths = append(proto.ProtoPaths, importedProto.ProtoPaths...)
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

func builtIn(dataType string) bool {
	for _, protoDataType := range protoDataTypes {
		if dataType == protoDataType {
			return true
		}
	}

	return false
}
