package execute

import (
	"io"
	"strconv"
	"strings"
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

func (p *Parser) parseKeyValues() (plan *Plan, err error) {
	token, _ := p.scan()
	if token == WORD {

	}
	token, _ = p.scan()
	if token == COLON {

	}
	token, _ = p.scan()
	if token == WORD {

	}

	token, _ = p.scan()
	if token == COMMA {
		p.parseKeyValues()
	}




	return nil, nil
}

func (p *Parser) parseBlockContent() (plan *Plan, err error) {
	token, _ := p.scan()
	if token == WORD {


	} else {
		// err
	}

	token, _ = p.scan()
	if token == COMMA {
		p.parseKeyValues()
	}
	return nil, nil
}

func (p *Parser) parseBlock() (plan *Plan, err error) {

	times := 1
	mode := "s"

	token, val := p.scan()

	if token == NUMBER {
		times, _ = strconv.Atoi(val)
		token, val = p.scan()
	} else {
		times = 1
	}
	if token == WORD {
		if val == "p" {

		} else if val == "s" {

		} else {
			return nil, nil
		}

		token, val = p.scan()
	}
	if token == BLOCK_START {
		token, val = p.scan()
		if isStartBlock(token) {
			p.unscan()
			_, _ = p.parseBlock()
		} else {
			_, _ = p.parseBlockContent()
		}
		token, val = p.scan()
	} else {
		// err
	}

	_ = times
	_ = mode
	return nil, nil

}

func isStartBlock(token Token) bool {
	if token == BLOCK_START || token == NUMBER {
		return true
	}
	return false
}

func Parse(molecule string) (plan *Plan, err error) {

	parser := NewParser(strings.NewReader(molecule))
	parser.parseBlock()
	return nil, nil
}
