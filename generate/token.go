package generate

const (
	TokenFileKey          = "fileKey"
	TokenFileVal          = "fileVal"
	TokenRPCName          = "rpcName"
	TokenRPCIn            = "rpcIn"
	TokenRPCOut           = "rpcOut"
	TokenFileOptionKey    = "fileOptionKey"
	TokenFileOptionVal    = "fileOptionVal"
	TokenServiceKey       = "serviceKey"
	TokenServiceOptionKey = "serviceOptionKey"
	TokenServiceOptionVal = "serviceOptionVal"
	TokenRPCOptionKey     = "rpcOptionKey"
	TokenRPCOptionVal     = "rpcOptionVal"
	TokenMessageKey       = "messageKey"
	TokenFieldDataType    = "fieldDataType"
	TokenFieldKey         = "fieldKey"
	TokenFieldOptionKey   = "fieldOptionKey"
	TokenFieldOptionVal   = "fieldOptionVal"
	TokenRPCDone          = "rpcDone"
	TokenMessageDone      = "messageDone"
	TokenServiceDone      = "serviceDone"
	TokenFieldDone        = "fieldDone"
	TokenFieldRule        = "fieldRule"
)

type Token struct {
	Name  string
	Value string
}
