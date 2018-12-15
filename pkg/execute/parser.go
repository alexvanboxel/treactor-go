package execute

import (
	"errors"
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

func (p *Parser) parseKeyValues(kv map[string]string) (out map[string]string, err error) {
	token, key := p.scan()
	if token != WORD {
		return nil, errors.New("KV need key")
	}
	token, _ = p.scan()
	if token != COLON {
		return nil, errors.New("KV need :")
	}
	token, value := p.scan()
	if token == WORD || token == NUMBER {
		kv[key] = value
	} else {
		return nil, errors.New("KV needs value")
	}

	token, _ = p.scan()
	if token == COMMA {
		return p.parseKeyValues(kv)
	} else {
		p.unscan()
	}
	return kv, nil
}

func (p *Parser) parseBlockContent() (plan Plan, err error) {
	token, _ := p.scan()
	if token == WORD {

	} else {
		return nil, errors.New("")
	}

	token, _ = p.scan()
	if token == COMMA {
		_, err := p.parseKeyValues(make(map[string]string))
		if err != nil {
			return nil, err
		}
		token, _ = p.scan()
	}

	if token != BLOCK_END {
		return nil, errors.New("Expected Block End")
	}
	return nil, nil
}

func (p *Parser) parseBlock() (plan Plan, err error) {

	times := 1
	mode := "s"
	var next Plan

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
			return nil, errors.New("Only s or p accepted")
		}

		token, val = p.scan()
	}
	if token == BLOCK_START {
		token, val = p.scan()
		if isStartBlock(token) {
			p.unscan()
			next, err = p.parseBlock()
			if err != nil {
				return nil, err
			}
		} else {
			p.unscan()
			next, err = p.parseBlockContent()
			if err != nil {
				return nil, err
			}
		}
		token, val = p.scan()
	} else {
		return nil, errors.New("Unknown token for block")
	}

	if times > 1 {
		return &Repeat{
			times: times,
			mode:  mode,
			block: next,
		}, nil
	}

	return nil, nil

}

func isStartBlock(token Token) bool {
	if token == BLOCK_START || token == NUMBER {
		return true
	}
	return false
}

func Parse(molecule string) (plan Plan, err error) {

	parser := NewParser(strings.NewReader(molecule))
	return parser.parseBlock()
}
