package main

import (
	"bytes"
	"fmt"
	"strings"
)

type token int

const (
	eofTok token = iota
	semiTok
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
	// for balancing parens and braces
	parens []rune
	insertSemi bool
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
	p.lex()
	return &p
}

func (p *parser) push(c rune) {
	p.parens = append(p.parens, c)
}

func (p *parser) pop(c rune) {
	if len(p.parens) <= 0 {
		p.errorf("unbalanced braces or parentheses\n")
	}
	top := p.top()
	expected := ')'
	if top == '{' {
		expected = '}'
	}
	if c != expected {
		p.errorf("expected %c\n", expected)
	}
	p.parens = p.parens[:len(p.parens)-1]
}

func (p *parser) top() rune {
	if len(p.parens) > 0 {
		return p.parens[len(p.parens)-1]
	}
	return 0
}

func (p *parser) errorf(f string, args ...interface{}) {
	panic(fmt.Sprintf(f, args...))
}

func (p *parser) expect(toks ...token) {
	for _, tok := range toks {
		if p.tok == tok {
			p.lex()
			return
		}
	}
	if len(toks) == 1 {
		p.errorf("expected %q, got %q\n", toks[0], p.tok)
	} else {
		var buf bytes.Buffer
		for _, tok := range toks[:len(toks)-1] {
			buf.WriteString(string(tok))
			buf.WriteString(", ")
		}
		buf.WriteString(" or ")
		buf.WriteString(string(toks[len(toks)-1]))
		p.errorf("expected %v, got %q\n", buf, p.tok)
	}
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
	insertSemi := p.insertSemi && p.top() != '('
	for p.ch == ' ' || p.ch == '\t' || !insertSemi && p.ch == '\n' {
		p.next()
	}
}

func (p *parser) readAtom() {
	for strings.IndexRune("\n \t#;${}", p.ch) < 0 {
		p.next()
	}
}

func (p *parser) readQuote() string {
	p.next()
	offset := p.offset
	for p.ch != '\'' {
		p.next()
	}
	lit := string(p.src[offset:p.offset])
	p.next()
	fmt.Println("lit is", lit)
	return lit
}

func special(c rune) bool {
	return strings.IndexRune(";{}()", c) >= 0
}

func (p *parser) lex() {
	p.skipSpace()

	offset := p.offset
	switch p.ch {
	case eof:
		if p.insertSemi {
			p.tok = semiTok
			p.insertSemi = false
			break
		}
		p.tok = eofTok
	case '\'':
		p.lit = p.readQuote()
		p.tok = atomTok
		return
	case '$':
		p.next()
		p.readAtom()
		p.tok = varTok
		p.insertSemi = true
	case '#':
		p.next()
		for p.ch != '\n' {
			p.next()
		}
		fallthrough
	case '\n', ';':
		p.next()
		p.tok = semiTok
		p.insertSemi = false
	default:
		if special(p.ch) {
			p.tok = token(p.ch)
			switch p.ch {
			case '{', '(':
				p.push(p.ch)
				p.insertSemi = false
			case '}', ')':
				p.pop(p.ch)
				fallthrough
			default:
				p.insertSemi = true
			}
			p.next()
			break
		}
		p.readAtom()
		p.tok = atomTok
		p.insertSemi = true
	}
	p.lit = string(p.src[offset:p.offset])
}

func (p *parser) parseParenExpr() []expr {
	p.lex()
	args := p.parseArgList()
	p.expect(')')
	return args
}

func (p *parser) parseBlock() block {
	p.lex()
	cmds := p.parseCommandList()
	p.expect('}')
	return block{cmds}
}

func (p *parser) parseArgList() []expr {
	var args []expr
	loop:
	for {
		switch p.tok {
		case atomTok:
			args = append(args, word(p.lit))
			p.lex()
		case varTok:
			args = append(args, localVar(p.lit))
			p.lex()
		case '(':
			args = append(args, p.parseParenExpr()...)
		case '{':
			 args = append(args, p.parseBlock())
		default:
			 break loop
		}
	}
	return args
}

func (p *parser) parseCommand() *completeCmd {
	var c cmd
	switch p.tok {
	case atomTok:
		c = word(p.lit)
		p.lex()
	case '{':
		c = p.parseBlock()
	default:
		p.expect('{', atomTok)
	}

	args := p.parseArgList()
	if p.tok != '}' {
		p.expect(semiTok)
	}
	return &completeCmd{c, args}
}

func (p *parser) parseCommandList() []*completeCmd {
	commands := make([]*completeCmd, 0, 1)
	for p.tok != eofTok && p.tok != '}' {
		commands = append(commands, p.parseCommand())
	}
	return commands
}
