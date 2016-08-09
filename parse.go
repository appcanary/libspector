package libspector

import "fmt"

type ParseError struct {
	Line    string
	Command string
}

func (err *ParseError) Error() string {
	if err.Command != "" {
		return fmt.Sprintf("failed to parse command [%s] output: %s", err.Command, err.Line)
	}
	return fmt.Sprintf("failed to parse: %s", err.Line)
}

func ErrParse(output string) error {
	return &ParseError{Line: output}
}

func IsParseError(err error) bool {
	_, ok := err.(*ParseError)
	return ok
}
