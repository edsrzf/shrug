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
	for _, c := range builtins {
		top["fn-" + c.name] = c
	}
	ctx.vars = []map[string]val{top}
	return &ctx
}

func (c *context) copy() *context {
	ctx := *c
	return &ctx
}

func (c *context) let(name string, v val) {
	c.vars[len(c.vars)-1][name] = v
}

func (c *context) set(name string, v val) {
	for i := len(c.vars) - 1; i >= 0; i-- {
		if v, ok := c.vars[i][name]; ok {
			c.vars[i][name] = v
			return
		}
	}
	c.vars[0][name] = v
}

func (c *context) lookupLocal(name string) val {
	for i := len(c.vars) - 1; i >= 0; i-- {
		if v, ok := c.vars[i][name]; ok {
			return v
		}
	}
	return nilVal{}
}

func (c *context) lookupFunc(name string) cmd {
	name = "fn-" + name
	for i := len(c.vars) - 1; i >= 0; i-- {
		if v, ok := c.vars[i][name]; ok {
			if f, ok := v.(cmd); ok {
				return f
			}
		}
	}
	return nil
}

func (c *context) pushScope() {
	c.vars = append(c.vars, map[string]val{})
}

func (c *context) popScope() {
	c.vars = c.vars[:len(c.vars)-1]
}

type builtinCmd struct {
	name string
	f    func(args []val, ctx *context) val
}

func (c *builtinCmd) exec(args []val, ctx *context) val {
	return c.f(args, ctx)
}

func (c *builtinCmd) eval(ctx *context) val {
	return c
}

func (c *builtinCmd) String() string { return "$&" + c.name }

func (c *builtinCmd) bool() bool { return true }

type block struct {
	cmds []*completeCmd
}

func (b block) exec(args []val, ctx *context) val {
	ctx.pushScope()
	defer ctx.popScope()

	var ret val = nilVal{}
	for _, cmd := range b.cmds {
		ret = cmd.exec(nil, ctx)
	}
	return ret
}

func (b block) eval(ctx *context) val { return b }

func (b block) String() string {
	// TODO
	return ""
}

func (b block) bool() bool { return true }

type argBlock struct {
	argNames []string
	cmds []*completeCmd
}

func (b argBlock) exec(args []val, ctx *context) val {
	ctx.pushScope()
	defer ctx.popScope()

	switch {
	case len(args) == len(b.argNames):
		for i, arg := range args {
			ctx.let(b.argNames[i], arg)
		}
	case len(args) > len(b.argNames):
		for i, argName := range b.argNames[:len(b.argNames)-1] {
			ctx.let(argName, args[i])
		}
		ctx.let(b.argNames[len(b.argNames)-1], argsToVal(args[len(b.argNames):]))
	}

	var ret val = nilVal{}
	for _, cmd := range b.cmds {
		ret = cmd.exec(nil, ctx)
	}
	return ret
}

func (b argBlock) eval(ctx *context) val { return b }

func (b argBlock) String() string {
	// TODO
	return ""
}

func (b argBlock) bool() bool { return true }

// A complete command that already has all its arguments. Its exec ignores
// the args parameter.
type completeCmd struct {
	cmd cmd
	args []expr
}

func (c *completeCmd) exec(args []val, ctx *context) val {
	termArgs := make([]val, 0, len(c.args))
	for _, expr := range c.args {
		val := expr.eval(ctx)
		if listVal, ok := val.(list); ok {
			termArgs = append(termArgs, listVal...)
		} else {
			termArgs = append(termArgs, val)
		}
	}
	ret := c.cmd.exec(termArgs, ctx)
	ctx.set("?", ret)
	return ret
}

type cmd interface {
	exec(args []val, ctx *context) val
}
