package main

import (
	"io"
)

type context struct {
	vars []map[string]termVal
	stdin io.Reader
	stdout, stderr io.Writer
}

func newCtx() *context {
	var ctx context
	top := map[string]termVal{}
	for _, c := range builtins {
		top["fn-" + c.name] = c
	}
	ctx.vars = []map[string]termVal{top}
	return &ctx
}

func (c *context) copy() *context {
	ctx := *c
	return &ctx
}

func (e *context) set(name string, v termVal) {
	e.vars[len(e.vars)-1][name] = v
}

func (e *context) lookupLocal(name string) termVal {
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

type builtinCmd struct {
	name string
	f    func(args []termVal, ctx *context) termVal
}

func (c *builtinCmd) exec(args []termVal, ctx *context) termVal {
	return c.f(args, ctx)
}

func (c *builtinCmd) eval(ctx *context) termVal {
	return c
}

func (c *builtinCmd) String() string { return "$&" + c.name }

func (c *builtinCmd) bool() bool { return true }

type lambda struct {
	cmds []*completeCmd
}

func (l lambda) exec(args []termVal, ctx *context) termVal {
	var ret termVal = nilVal{}
	for _, cmd := range l.cmds {
		ret = cmd.exec(nil, ctx)
	}
	return ret
}

func (l lambda) eval(ctx *context) termVal {
	return l
}

func (l lambda) String() string {
	// TODO
	return ""
}

func (l lambda) bool() bool { return true }

// A complete command that already has all its arguments. Its exec ignores
// the args parameter.
type completeCmd struct {
	cmd cmd
	args []val
}

func (c *completeCmd) exec(args []termVal, ctx *context) termVal {
	termArgs := make([]termVal, 0, len(c.args))
	for _, val := range c.args {
		termArgs = append(termArgs, val.eval(ctx))
	}
	return c.cmd.exec(termArgs, ctx)
}

type cmd interface {
	exec(args []termVal, ctx *context) termVal
}
