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
	{"set var hi; echo $var", "hi\n"},
	{"{echo hi}", "hi\n"},
	{"{\necho hi}", "hi\n"},
	{"if {test 0} {echo true} {echo false}", "true\n"},
	{"if {test 0} {echo true}", "true\n"},
	{"if {test} {echo true} {echo false}", "false\n"},
	{"if {test} {echo true}", ""},
	{"for var {echo hi}", ""},
	{"for letter a b c {echo $letter}", "a\nb\nc\n"},
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
			cmd.exec(nil, ctx)
		}
		if output := buf.String(); output != test.output {
			t.Errorf("expected\n%q\ngot\n%q", test.output, output)
		}
	}
}
