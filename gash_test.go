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
	{"echo $hi", "\n"},
}

func TestCommand(t *testing.T) {
	var buf bytes.Buffer
	for _, test := range basicTests {
		buf.Reset()
		p := newParser(test.cmd + "\n")
		cmds := p.parseCommandList()
		cmds[0].exec(nil, &buf, nil, &environment{})
		if output := buf.String(); output != test.output {
			t.Errorf("expected\n%q\ngot\n%q", test.output, output)
		}
	}
}
