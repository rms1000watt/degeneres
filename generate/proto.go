package generate

// FileState
type Proto struct {
	Syntax   string    // FileSyntaxState
	Package  string    // FilePackageState
	Imports  []string  // FileImportsState
	Options  []Option  // FileOptionsState
	Services []Service // FileServicesState
	Messages []Message // FileMessagesState
}

type Option struct {
	Name  string
	Value string
}

type Service struct {
	Name string
	RPCs []RPC
}

type RPC struct {
	Name    string
	Input   string
	Output  string
	Options []Option
}

type Message struct {
	Name    string
	Options []Option
	Fields  []Field
}

type Field struct {
	Name     string
	Type     string
	Position string
	Options  []Option
}
