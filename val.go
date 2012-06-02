package main

import (
	"fmt"
	"math/big"
	"os/exec"
	"syscall"
)

type val interface {
	eval(ctx *context) val
	String() string
}

type nilVal struct{}

func (v nilVal) eval(ctx *context) val {
	return v
}

func (v nilVal) String() string { return "" }

type localVar string

func (v localVar) eval(ctx *context) val { return ctx.lookupLocal(string(v[1:])) }

func (v localVar) String() string {
	panic("shouldn't be called")
}

type word string

func (w word) exec(args []val, ctx *context) int {
	if f := ctx.lookupFunc(string(w)); f != nil {
		return f.exec(args, ctx)
	}
	path, err := exec.LookPath(string(w))
	if err != nil {
		fmt.Printf("%s: command not found\n", w)
		return 127
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
			return ee.Sys().(*syscall.WaitStatus).ExitStatus()
		}
		panic(err)
	}
	return 0
}

func (w word) eval(ctx *context) val { return w }

func (w word) String() string { return string(w) }

type integer struct {
	val *big.Int
}

type rational struct {
	val *big.Rat
}

type block struct {
}

type list struct {
}
