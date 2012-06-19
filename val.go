package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"syscall"
)

type val interface {
	eval(ctx *context) termVal
}

type termVal interface {
	val
	String() string
	bool() bool
}

type nilVal struct{}

func (v nilVal) eval(ctx *context) termVal {
	return v
}

func (v nilVal) String() string { return "" }

func (v nilVal) bool() bool { return true }

type localVar string

func (v localVar) eval(ctx *context) termVal { return ctx.lookupLocal(string(v[1:])) }

type word string

func (w word) exec(args []termVal, ctx *context) termVal {
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
		strArgs[i] = arg.String()
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

func (w word) eval(ctx *context) termVal { return w }

func (w word) String() string { return string(w) }

func (w word) bool() bool { return w == "" || w == "0" }

type intVal int

func (v intVal) eval(ctx *context) termVal { return v }

func (v intVal) String() string { return strconv.Itoa(int(v)) }

func (v intVal) bool() bool { return v == 0 }

type list []termVal

func (v list) eval(ctx *context) termVal { return v }

func (v list) String() string {
	var buf bytes.Buffer
	for i, el := range v {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(el.String())
	}
	return buf.String()
}

func (v list) bool() bool {
	for _, el := range v {
		if !el.bool() {
			return false
		}
	}
	return true
}
