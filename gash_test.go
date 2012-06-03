package main

import (
	"bytes"
	"testing"
)

var basicTests = []struct{
	cmd string
	output string
}{
	{"echo hi", "hi\n"},
	{"echo hi; echo hi", "hi\nhi\n"},
	{"echo $hi", "\n"},
	{"set var hi", ""},
	{"{echo hi}", "hi\n"},
}

func TestCommand(t *testing.T) {
	var buf bytes.Buffer
	for _, test := range basicTests {
		buf.Reset()
		p := newParser(test.cmd + "\n")
		cmds := p.parseCommandList()
		ctx := newCtx()
		ctx.stdout = &buf
		for _, cmd := range cmds {
			cmd.exec(ctx)
		}
		if output := buf.String(); output != test.output {
			t.Errorf("expected\n%q\ngot\n%q", test.output, output)
		}
	}
}
