package evaluator

import (
	"Monkey/ast"
	"Monkey/object"
)

// DefineMacros notes down the macro of the ast tree and remove them
func DefineMacros(program *ast.Program, env *object.Environment) {
	var definitions []int

	// Find macro definitions
	for i, statement := range program.Statements {
		if isMacroDefinition(statement) {
			addMacro(statement, env)
			definitions = append(definitions, i)
		}
	}

	// Remove them from the ast tree
	for i := len(definitions) - 1; i >= 0; i -= 1 {
		definitionIndex := definitions[i]
		// splice at index
		program.Statements = append(
			program.Statements[:definitionIndex],
			program.Statements[definitionIndex+1:]...,
		)
	}
}

// isMacroDefinition return where the node is a macro definition
func isMacroDefinition(node ast.Statement) bool {
	// is statement
	letStatement, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}

	// does it contain a macro literal
	_, ok = letStatement.Value.(*ast.MacroLiteral)
	return ok
}

// addMacro adds a macro to the environment
func addMacro(stmt ast.Statement, env *object.Environment) {
	// ignore error as it is checked in isMacroDefinition
	letStatement, _ := stmt.(*ast.LetStatement)
	macroLiteral, _ := letStatement.Value.(*ast.MacroLiteral)

	// create macro object
	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Env:        env,
		Body:       macroLiteral.Body,
	}

	// store to environment
	env.Store(letStatement.Name.Value, macro)
}

// ExpandMacros takes a root program node and expands all the macros calls with the macros defined
// in the environment passed in
func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExpression, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		// check for macro call
		macro, ok := isMacroCall(callExpression, env)
		if !ok {
			return node
		}

		// get macro args and expand the environment
		args := quoteArgs(callExpression)
		evalEnv := extendMacroEnv(macro, args)

		// eval macro body
		evaluated := Eval(macro.Body, evalEnv)

		// return the new ast node
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			// TODO Change this panic call to an error
			panic("Only support returning ast node from macros")
		}

		return quote.Node
	})
}

/// isMacroCall returns whether the callExpression is a macro call
func isMacroCall(
	exp *ast.CallExpression,
	env *object.Environment,
) (*object.Macro, bool) {
	identifier, ok := exp.Function.(*ast.Identifier)
	if !ok {
		return nil, false
	}

	obj, ok := env.Get(identifier.Value)
	if !ok {
		return nil, false
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}

	return macro, true
}

// quoteArgs returns a list of quotes from a ast.CallExpression
func quoteArgs(exp *ast.CallExpression) []*object.Quote {
	var args []*object.Quote

	for _, a := range exp.Arguments {
		args = append(args, &object.Quote{Node: a})
	}

	return args
}

// extendMacroEnv returns the extended environment for a macro call
func extendMacroEnv(
	macro *object.Macro,
	args []*object.Quote,
) *object.Environment {
	extended := object.NewEnclosingEnvironment(macro.Env)

	for paramIdx, param := range macro.Parameters {
		extended.Store(param.Value, args[paramIdx])
	}

	return extended
}
