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

type lambda struct {
	cmds []*completeCmd
}

func (l lambda) exec(args []val, ctx *context) int {
	ret := 0
	for _, cmd := range l.cmds {
		ret = cmd.exec(nil, ctx)
	}
	return ret
}

func (l lambda) eval(ctx *context) val {
	return l
}

func (l lambda) String() string {
	// TODO
	return ""
}

// A complete command that already has all its arguments. Its exec ignores
// the args parameter.
type completeCmd struct {
	cmd cmd
	args []val
}

func (c *completeCmd) exec(args []val, ctx *context) int {
	return c.cmd.exec(c.args, ctx)
}

type cmd interface {
	exec(args []val, ctx *context) int
}
