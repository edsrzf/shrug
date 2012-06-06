package main

import (
	"bytes"
	"testing"
)

var basicTests = []struct{
	cmd string
	output string
	ret val
}{
	{"echo hi", "hi\n", intVal(0)},
	{"echo hi; echo hi", "hi\nhi\n", intVal(0)},
	{"echo $hi", "\n", intVal(0)},
	{"set var hi", "", nilVal{}},
	{"set var hi; echo $var", "hi\n", intVal(0)},
	{"{echo hi}", "hi\n", intVal(0)},
	{"{\necho hi}", "hi\n", intVal(0)},
	{"if {test 0} {echo true} {echo false}", "true\n", intVal(0)},
	{"if {test 0} {echo true}", "true\n", intVal(0)},
	{"if {test} {echo true} {echo false}", "false\n", intVal(0)},
	{"if {test} {echo true}", "", intVal(0)},
	{"for var {echo hi}", "", intVal(0)},
	{"for letter a b c {echo $letter}", "a\nb\nc\n", intVal(0)},
	{"result", "", nilVal{}},
	{"result 1", "", word("1")},
	{"result 1 2 3", "", list{word("1"), word("2"), word("3")}},
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
