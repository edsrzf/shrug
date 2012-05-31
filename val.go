package main

import (
	"fmt"
	"io"
	"math/big"
	"os/exec"
	"syscall"
)

type val interface {
	eval(env *environment) val
	String() string
}

type word string

func (w word) exec(args []val, stdin io.Reader, stdout, stderr io.Writer, env *environment) int {
	path, err := exec.LookPath(string(w))
	if err != nil {
		fmt.Printf("%s: command not found\n", w)
		return 127
	}
	strArgs := make([]string, len(args))
	for i, arg := range args {
		strArgs[i] = arg.String()
	}
	cmd := exec.Command(path, strArgs...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
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

func (w word) eval(env *environment) val { return w }

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
