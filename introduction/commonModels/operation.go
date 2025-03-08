package commonModels

// Operation describes type signature of an execution
// its Signature is a character symbol of the operation displayed in the console
type Operation struct {
	Execute   func(args ...any) any
	Signature byte // only a getter function should be exported
}
