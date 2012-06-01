package main

import (
	"fmt"
	"strings"
)

type token int

const (
	eofTok token = iota
	lfTok
	atomTok
	varTok
	countTok
	evalTok // `{
)

type lexMode int

const (
	normalMode lexMode = iota
	dquoteMode
)

type parser struct {
	ch rune
	tok token
	lit string
	mode lexMode
	src []byte
	offset int
}

const eof = -1

func newParser(src string) *parser {
	var p parser
	p.src = []byte(src)
	if len(src) > 0 {
		p.ch = rune(src[0])
	} else {
		p.ch = eof
	}
	p.offset = 0
	return &p
}

func (p *parser) errorf(f string, args ...interface{}) {
	panic(fmt.Sprintf(f, args...))
}

func (p *parser) next() {
	if p.offset < len(p.src) - 1 {
		// TODO: unicode
		p.offset++
		p.ch = rune(p.src[p.offset])
	} else {
		p.ch = eof
	}
}

func (p *parser) skipSpace() {
	for p.ch == ' ' || p.ch == '\t' {
		p.next()
	}
}

func (p *parser) readAtom() {
	for strings.IndexRune("\n \t#;${}", p.ch) < 0 {
		p.next()
	}
}

func special(c rune) bool {
	return strings.IndexRune(";{}", c) >= 0
}

func (p *parser) lex() {
	p.skipSpace()

	offset := p.offset
	switch p.ch {
	case eof:
		p.tok = eofTok
	case '$':
		p.next()
		p.readAtom()
		p.tok = varTok
	case '#':
		p.next()
		for p.ch != '\n' {
			p.next()
		}
		fallthrough
	case '\n':
		p.next()
		p.tok = lfTok
	default:
		if special(p.ch) {
			p.next()
			p.tok = token(p.ch)
			break
		}
		p.readAtom()
		p.tok = atomTok
	}
	p.lit = string(p.src[offset:p.offset])
}

func (p *parser) parseCommand() *command {
	if p.tok != atomTok {
		p.errorf("expected atom, got %d\n", p.tok)
	}
	w := word(p.lit)
	p.lex()
	var args []val
	for p.tok != lfTok {
		switch p.tok {
		case atomTok:
			args = append(args, word(p.lit))
		case varTok:
			args = append(args, localVar(p.lit))
		}
		p.lex()
	}
	p.lex()
	return &command{w, args}
}

func (p *parser) parseCommandList() []*command {
	p.lex()
	commands := make([]*command, 0, 1)
	for p.tok != eofTok {
		commands = append(commands, p.parseCommand())
	}
	return commands
}
