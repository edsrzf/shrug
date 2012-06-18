package main

import (
	"fmt"
	"io"
	"os"

	"github.com/edsrzf/fineline"
)

func main() {
	ctx := newCtx()
	ctx.stdin = os.Stdin
	ctx.stdout = os.Stdout
	ctx.stderr = os.Stderr

	l := fineline.NewLineReader()
	l.Prompt = "$ "
	l.SetMaxHistory(10)

	for {
		str, err := l.Read(nil)
		if err != nil {
			if err != io.EOF {
				fmt.Println("error", err)
			} else {
				fmt.Println()
			}
			break
		}
		parser := newParser(str + "\n")
		cmd := parser.parseCommand()
		cmd.exec(nil, ctx)
	}
}
