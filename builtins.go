package main

import (
	"fmt"
	"os"
)

var builtins = []*builtinCmd{
	{"and", andCmd},
	{"cd", cdCmd},
	{"create", createCmd},
	{"for", forCmd},
	{"if", ifCmd},
	{"let", letCmd},
	{"pipe", pipeCmd},
	{"result", resultCmd},
	{"set", setCmd},
}

func andCmd(args []val, ctx *context) val {
	var ret val = nilVal{}
	for _, arg := range args {
		c, ok := arg.(cmd)
		if !ok {
		}
		if ret = c.exec(nil, ctx); !ret.bool() {
			break
		}
	}
	return ret
}

func cdCmd(args []val, ctx *context) val {
	var dir string
	if len(args) == 0 {
		dir = os.Getenv("HOME")
	} else {
		dir = args[0].String()
	}
	err := os.Chdir(dir)
	if err != nil {
		fmt.Fprintln(ctx.stderr, "cd:", err)
		return intVal(1)
	}
	return intVal(0)
}

func createCmd(args []val, ctx *context) val {
	usage := func() val {
		ctx.stderr.Write([]byte("create: usage: create fd filename cmd"))
		return intVal(1)
	}
	if len(args) != 3 {
		return usage()
	}
	cmdArg, ok := args[2].(cmd)
	if !ok {
		return usage()
	}
	file, err := os.Create(args[1].String())
	if err != nil {
		ctx.stderr.Write([]byte("create: unable to create file"))
		return intVal(1)
	}
	defer file.Close()

	switch args[0].String() {
	case "1":
		oldStdout := ctx.stdout
		ctx.stdout = file
		defer func() {
			ctx.stdout = oldStdout
		}()
	case "2":
		oldStderr := ctx.stdout
		ctx.stderr = file
		defer func() {
			ctx.stdout = oldStderr
		}()
	default:
		ctx.stderr.Write([]byte("create: can only redirect fds 1 and 2"))
		return intVal(1)
	}

	return cmdArg.exec(nil, ctx)
}

func ifCmd(args []val, ctx *context) val {
	if len(args) < 1 {
		ctx.stderr.Write([]byte("if: usage: if cond [ iftrue ] [ iffalse ]"))
		return intVal(1)
	}
	cond, ok := args[0].(cmd)
	if !ok {
	}
	ret := cond.exec(nil, ctx)
	if ret.bool() {
		if len(args) >= 2 {
			iftrue, ok := args[1].(cmd)
			if ok {
				iftrue.exec(nil, ctx)
			}
		}
	} else if len(args) >= 3 {
		iffalse, ok := args[2].(cmd)
		if ok {
			iffalse.exec(nil, ctx)
		}
	}
	return nilVal{}
}

func forCmd(args []val, ctx *context) val {
	if len(args) < 2 {
		ctx.stderr.Write([]byte("for: usage: for variable [ list ... ] body"))
		return intVal(1)
	}
	varname, ok := args[0].(word)
	if !ok {
		ctx.stderr.Write([]byte("for: invalid variable name"))
		return intVal(1)
	}
	body, ok := args[len(args)-1].(cmd)
	if !ok {
		ctx.stderr.Write([]byte("for: invalid body"))
		return intVal(1)
	}
	for _, item := range args[1 : len(args)-1] {
		ctx.set(string(varname), item)
		body.exec(nil, ctx)
	}
	return nilVal{}
}

func letCmd(args []val, ctx *context) val {
	if len(args) < 2 {
		ctx.stderr.Write([]byte("let: usage: let variable value"))
		return intVal(1)
	}
	varname, ok := args[0].(word)
	if !ok {
		ctx.stderr.Write([]byte("let: invalid variable name"))
		return intVal(1)
	}
	ctx.let(string(varname), argsToVal(args[1:]))
	return nilVal{}
}

func pipeCmd(args []val, ctx *context) val {
	stdin := ctx.stdin
	if len(args) == 0 {
		return nilVal{}
	}
	if len(args) > 1 {
		for _, arg := range args[:len(args)-1] {
			pipeR, pipeW, err := os.Pipe()
			if err != nil {
				continue
			}
			cmdArg, ok := arg.(cmd)
			if !ok {
				continue
			}
			subCtx := ctx.copy()
			subCtx.stdin = stdin
			subCtx.stdout = pipeW
			go func() {
				cmdArg.exec(nil, subCtx)
				pipeW.Close()
			}()
			stdin = pipeR
		}
	}
	cmdArg, ok := args[len(args)-1].(cmd)
	if !ok {
		return nilVal{}
	}
	oldStdin := ctx.stdin
	ctx.stdin = stdin
	defer func() {
		ctx.stdin = oldStdin
	}()
	return cmdArg.exec(nil, ctx)
}

func argsToVal(args []val) val {
	switch len(args) {
	case 0:
		return nilVal{}
	case 1:
		return args[0]
	}
	return list(args)
}

func resultCmd(args []val, ctx *context) val {
	return argsToVal(args)
}

func setCmd(args []val, ctx *context) val {
	if len(args) < 2 {
		ctx.stderr.Write([]byte("set: usage: set variable value"))
		return intVal(1)
	}
	varname, ok := args[0].(word)
	if !ok {
		ctx.stderr.Write([]byte("set: invalid variable name"))
		return intVal(1)
	}
	ctx.set(string(varname), argsToVal(args[1:]))
	return nilVal{}
}
