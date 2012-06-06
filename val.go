package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"syscall"
)

type val interface {
	eval(ctx *context) val
	bool() bool
	String() string
}

type nilVal struct{}

func (v nilVal) eval(ctx *context) val {
	return v
}

func (v nilVal) String() string { return "" }

func (v nilVal) bool() bool { return true }

type localVar string

func (v localVar) eval(ctx *context) val { return ctx.lookupLocal(string(v[1:])) }

func (v localVar) String() string {
	panic("shouldn't be called")
}

func (v localVar) bool() bool {
	panic("shouldn't be called")
}

type word string

func (w word) exec(args []val, ctx *context) val {
	if f := ctx.lookupFunc(string(w)); f != nil {
		return f.exec(args, ctx)
	}
	path, err := exec.LookPath(string(w))
	if err != nil {
		fmt.Printf("%s: command not found\n", w)
		return intVal(127)
	}
	strArgs := make([]string, len(args))
	for i, arg := range args {
		strArgs[i] = arg.eval(ctx).String()
	}
	cmd := exec.Command(path, strArgs...)
	cmd.Stdin = ctx.stdin
	cmd.Stdout = ctx.stdout
	cmd.Stderr = ctx.stderr
	err = cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			// TODO: make portable to Plan 9
			return intVal(ee.Sys().(syscall.WaitStatus).ExitStatus())
		}
		panic(err)
	}
	return nilVal{}
}

func (w word) eval(ctx *context) val { return w }

func (w word) String() string { return string(w) }

func (w word) bool() bool { return w == "" || w == "0" }

type intVal int

func (v intVal) eval(ctx *context) val { return v }

func (v intVal) String() string { return strconv.Itoa(int(v)) }

func (v intVal) bool() bool { return v == 0 }

type list struct {
}
