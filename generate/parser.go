package generate

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func Parse(tokens chan Token) (proto Proto) {
	log.Debug("Staring parser")
	defer log.Debug("Parser done")

	proto = NewProto()

	kv := KV{}
	rpc := RPC{}
	field := Field{}
	message := Message{}
	service := Service{}
	rpcOption := Option{}
	fileOption := Option{}
	fieldOption := Option{}
	serviceOption := Option{}

	for token := range tokens {
		switch token.Name {
		case TokenFileKey:
			kv.Key = token.Value
		case TokenFileVal:
			kv.Val = token.Value
			handleFileKV(kv, &proto)
			kv = KV{}
		case TokenFileOptionKey:
			fileOption.Name = token.Value
		case TokenFileOptionVal:
			fileOption.Value = token.Value
			proto.Options = append(proto.Options, fileOption)
			fileOption = Option{}
		case TokenServiceKey:
			service.Name = token.Value
		case TokenServiceOptionKey:
			serviceOption.Name = token.Value
		case TokenServiceOptionVal:
			serviceOption.Value = token.Value
			service.Options = append(service.Options, serviceOption)
			serviceOption = Option{}
		case TokenServiceDone:
			proto.Services = append(proto.Services, service)
			service = Service{}
		case TokenRPCName:
			rpc.Name = token.Value
		case TokenRPCIn:
			rpc.Input = token.Value
		case TokenRPCOut:
			rpc.Output = token.Value
		case TokenRPCOptionKey:
			rpcOption.Name = token.Value
		case TokenRPCOptionVal:
			rpcOption.Value = token.Value
			rpc.Options = append(rpc.Options, rpcOption)
			rpcOption = Option{}
		case TokenRPCDone:
			service.RPCs = append(service.RPCs, rpc)
			rpc = RPC{}
		case TokenMessageKey:
			message.Name = token.Value
		case TokenMessageDone:
			proto.Messages = append(proto.Messages, message)
			message = Message{}
		case TokenFieldDataType:
			field.DataType = token.Value
		case TokenFieldKey:
			field.Name = token.Value
		case TokenFieldDone:
			message.Fields = append(message.Fields, field)
			field = Field{}
		case TokenFieldOptionKey:
			fieldOption.Name = token.Value
		case TokenFieldOptionVal:
			fieldOption.Value = token.Value
			field.Options = append(field.Options, fieldOption)
			fieldOption = Option{}
		case TokenFieldRule:
			field.Rule = token.Value
		}
	}

	return
}

func handleFileKV(kv KV, proto *Proto) {
	switch strings.ToLower(kv.Key) {
	case FilePackage:
		proto.Package = kv.Val
	case FileSyntax:
		proto.Syntax = kv.Val
	case FileImport:
		proto.Imports = append(proto.Imports, kv.Val)
	}
}
