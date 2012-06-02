package main

import (
	"io"
)

type context struct {
	vars []map[string]val
	stdin io.Reader
	stdout, stderr io.Writer
}

func newCtx() *context {
	var ctx context
	top := map[string]val{}
	for name, f := range builtins {
		top["fn-" + name] = f
	}
	ctx.vars = []map[string]val{top}
	return &ctx
}

func (e *context) set(name string, v val) {
	e.vars[len(e.vars)-1][name] = v
}

func (e *context) lookupLocal(name string) val {
	for i := len(e.vars) - 1; i >= 0; i-- {
		if v, ok := e.vars[i][name]; ok {
			return v
		}
	}
	return nilVal{}
}

func (e *context) lookupFunc(name string) cmd {
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

type builtinCmd func(args []val, ctx *context) int

func (c builtinCmd) exec(args []val, ctx *context) int {
	return c(args, ctx)
}

func (c builtinCmd) eval(ctx *context) val {
	return nilVal{}
}

func (c builtinCmd) String() string {
	return ""
}

type function struct {
}

func (f *function) call(ctx *context, args ...val) int {
	return 0
}

func (f *function) eval(ctx *context) val {
	return nilVal{}
}

func (f *function) String() string {
	return ""
}

type command struct {
	cmd cmd
	args []val
}

func (c *command) exec(ctx *context) int {
	return c.cmd.exec(c.args, ctx)
}

type cmd interface {
	exec(args []val, ctx *context) int
}
