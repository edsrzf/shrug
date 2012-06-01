package main

import (
	"io"
)

type environment struct {
	vars []map[string]val
}

func (e *environment) lookup(name string) val {
	i := len(e.vars) - 1
	for i >= 0 {
		if v, ok := e.vars[i][name]; ok {
			return v
		}
	}
	return nilVal{}
}

func (e *environment) lookupFunc(name string) *function {
	return nil
}

type function struct {
}

func (f *function) call(env *environment, args ...val) int {
	return 0
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
