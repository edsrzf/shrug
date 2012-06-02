package main

import (
	"io"
)

type environment struct {
	vars []map[string]val
}

func newEnv() *environment {
	var env environment
	top := map[string]val{}
	for name, f := range builtins {
		top["fn-" + name] = f
	}
	env.vars = []map[string]val{top}
	return &env
}

func (e *environment) set(name string, v val) {
	e.vars[len(e.vars)-1][name] = v
}

func (e *environment) lookupLocal(name string) val {
	for i := len(e.vars) - 1; i >= 0; i-- {
		if v, ok := e.vars[i][name]; ok {
			return v
		}
	}
	return nilVal{}
}

func (e *environment) lookupFunc(name string) cmd {
	name = "fn-" + name
	for i := len(e.vars) - 1; i >= 0; i-- {
		if v, ok := e.vars[i][name]; ok {
			if f, ok := v.(cmd); ok {
				return f
			}
		}
	}
	return nil
}

type builtinCmd func(args []val, stdin io.Reader, stdout, stderr io.Writer, env *environment) int

func (c builtinCmd) exec(args []val, stdin io.Reader, stdout, stderr io.Writer, env *environment) int {
	return c(args, stdin, stdout, stderr, env)
}

func (c builtinCmd) eval(env *environment) val {
	return nilVal{}
}

func (c builtinCmd) String() string {
	return ""
}

type function struct {
}

func (f *function) call(env *environment, args ...val) int {
	return 0
}

func (f *function) eval(env *environment) val {
	return nilVal{}
}

func (f *function) String() string {
	return ""
}

type command struct {
	cmd cmd
	args []val
}

func (c *command) exec(stdin io.Reader, stdout, stderr io.Writer, env *environment) int {
	return c.cmd.exec(c.args, stdin, stdout, stderr, env)
}

type cmd interface {
	exec(args []val, stdin io.Reader, stdout, stderr io.Writer, env *environment) int
}
