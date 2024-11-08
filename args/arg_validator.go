package args

import "fmt"

func AssertArgsPresent(args map[string]Arg, nameList []string) error {
	for _, name := range nameList {
		_, ok := args[name]
		if !ok {
			return fmt.Errorf("required arg '%s' was not present in parsed args", name)
		}
	}

	return nil
}
