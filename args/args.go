package args

import (
	"errors"
	"os"
	"strings"
)

var ErrorNoArgs = errors.New("no arguments provided")
var ErrorInvalidArg = errors.New("invalid arg, must specify argument as name=\"value\"")

func parseArg(argString string) (*Arg, error) {
	const seperator = '='
	const delimiter = ' '

	argNameBuf := strings.Builder{}
	argValueBuf := strings.Builder{}

	parsingValue := false
	for _, argChar := range argString {

		if argChar == delimiter {
			break
		}

		if argChar == seperator {
			parsingValue = true
			continue
		}

		if parsingValue {
			argValueBuf.WriteRune(argChar)
		} else {
			argNameBuf.WriteRune(argChar)
		}
	}

	parsedArgName := argNameBuf.String()
	parsedArgValue := argValueBuf.String()

	// Strip the " off of each end
	parsedArgValue = parsedArgValue[1 : len(parsedArgValue)-1]

	if parsedArgName == "" || parsedArgValue == "" {
		return nil, ErrorInvalidArg
	}

	arg := &Arg{Name: parsedArgName, Value: parsedArgValue}

	return arg, nil
}

// The worlds worst arg parser
// Returns map that maps from arg_name -> Arg object
func ParseArgs() (map[string]Arg, error) {
	argStrings := os.Args[1:] // Skip the first arg (usually the path to the program)
	argLen := len(argStrings)

	if argLen == 0 {
		return nil, ErrorNoArgs
	}

	args := make(map[string]Arg, 0)
	for _, argString := range argStrings {
		arg, err := parseArg(argString)
		if err != nil {
			return nil, err
		}

		args[arg.Name] = *arg
	}

	return args, nil
}
