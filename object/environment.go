package object

// To store variables and functions alike
type Environment struct {
	store map[string]Object
	outer *Environment
}

// Create a new environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosingEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get an item from the environment
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	// Recursive loop to get the variable
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) Replace(name string, val Object) (Object, bool) {

	_, ok := e.store[name]
	if ok {
		e.store[name] = val
		return val, ok
	} else if e.outer != nil {
		_, ok = e.outer.Replace(name, val)
	}

	return val, ok
}

// Store an item to the environment
func (e *Environment) Store(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Delete(name string) bool {
	if _, ok := e.store[name]; ok {
		delete(e.store, name)
		return true
	}
	return false
}
