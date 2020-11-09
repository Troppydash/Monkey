package object

import "fmt"

type Argument struct {
	OType                  []ObjectType
	SuccessParsingFunction func(optional bool, arg Object) ParsedArgument
	Optional               bool
	Any                    bool
	VarArg                 bool
}

func NewArgument(function func(optional bool, arg Object) ParsedArgument,
	optional bool,
	any bool,
	vararg bool,
	types ...ObjectType) Argument {
	return Argument{
		SuccessParsingFunction: function,
		OType:                  types,
		Optional:               optional,
		Any:                    any,
		VarArg:                 vararg,
	}
}

func NewAnyVarargsArgument(
	function func(optional bool, arg Object) ParsedArgument) Argument {
	return Argument{
		SuccessParsingFunction: function,
		Any:                    true,
		VarArg:                 true,
	}
}

func NewOptionalArgument(
	function func(optional bool, arg Object) ParsedArgument,
	types ...ObjectType) Argument {
	return Argument{
		SuccessParsingFunction: function,
		Optional:               true,
		OType:                  types,
	}
}

func NewAnyOptionalArgument(
	function func(optional bool, arg Object) ParsedArgument,
	types ...ObjectType) Argument {
	return Argument{
		SuccessParsingFunction: function,
		Optional:               true,
		Any:                    true,
	}
}

type ArgumentsParser struct {
	Types []Argument
}

type ParsedArgument struct {
	Value interface{}
}

func NewParser(types ...Argument) *ArgumentsParser {
	return &ArgumentsParser{
		Types: types,
	}
}

func (ap *ArgumentsParser) AddArgument(newArgument Argument) {
	ap.Types = append(ap.Types, newArgument)
}

type ParseError struct {
	message string
}

func (pe *ParseError) Error() string {
	return pe.message
}

func (ap *ArgumentsParser) Parse(args []Object) ([]ParsedArgument, error) {
	var parsedArguments []ParsedArgument
	var lastArgument *Argument = nil
	lastIndex := 0
	for index, argType := range ap.Types {
		if argType.VarArg {
			lastArgument = &argType
			lastIndex = index
			break
		}
		if index >= len(args) {
			if argType.Optional {
				pArgu := argType.SuccessParsingFunction(true, nil)
				parsedArguments = append(parsedArguments, pArgu)
				break
			}
			return []ParsedArgument{}, NewParseError("not enough argument passed in. expected=%d, got=%d", len(ap.Types), len(args))
		}
		currentArgument := args[index]

		found := false
		if argType.Any {
			found = true
		} else {
			for _, v := range argType.OType {
				if v == currentArgument.Type() {
					found = true
				}
			}
		}

		if found {
			pArgu := argType.SuccessParsingFunction(false, currentArgument)
			parsedArguments = append(parsedArguments, pArgu)
		} else {
			return []ParsedArgument{}, NewParseError("type mismatch. expected=%q, got=%s", argType.OType, currentArgument.Type())
		}
	}
	if lastArgument != nil {
		argsLeft := args[lastIndex:]
		for _, arg := range argsLeft {
			found := false
			if lastArgument.Any {
				found = true
			} else {
				for _, v := range lastArgument.OType {
					if v == arg.Type() {
						found = true
					}
				}
			}

			if found {
				pArgu := lastArgument.SuccessParsingFunction(false, arg)
				parsedArguments = append(parsedArguments, pArgu)
			} else {
				return []ParsedArgument{}, NewParseError("type mismatch. expected=%q, got=%s", lastArgument.OType, arg.Type())
			}
		}
	}

	return parsedArguments, nil
}

func NewParseError(str string, format ...interface{}) *ParseError {
	return &ParseError{
		message: fmt.Sprintf(str, format...),
	}
}
