package parser

import (
	"fmt"
	"strconv"

	"github.com/NeowayLabs/abad/ast"
	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/lexer"
	"github.com/NeowayLabs/abad/token"
)

type (
	Parser struct {
		tokens    <-chan lexer.Tokval
		lookahead []lexer.Tokval

		filename string
	}

	parserfn func(*Parser) (ast.Node, error)
)

// used when the tokens is over
var tokEOF = lexer.Tokval{
	Type: token.EOF,
}

// literal parsers
var litParsers = map[token.Type]parserfn{
	token.Illegal:     parseIllegal,
	token.Decimal:     parseDecimal,
	token.Hexadecimal: parseHex,
}

// Parse input source into an AST representation.
func Parse(fname string, code string) (*ast.Program, error) {
	p := Parser{
		tokens:   lexer.Lex(utf16.Encode(code)),
		filename: fname,
	}

	return p.parse()
}

func (p *Parser) parse() (*ast.Program, error) {
	var nodes []ast.Node

	for {
		p.scry(1)

		tok := p.lookahead[0]
		if tok.Type == token.EOF {
			break
		}

		parser, ok := litParsers[tok.Type]
		if !ok {
			return nil, p.errorf(tok, "not implemented: %s", tok)
		}

		node, err := parser(p)
		if err != nil {
			return nil, err
		}

		// parsers should not leave tokens not processed
		// in the lookahead buffer.
		if len(p.lookahead) != 0 {
			panic("parsers not handling lookahead correctly")
		}

		nodes = append(nodes, node)
	}

	return &ast.Program{
		Nodes: nodes,
	}, nil
}

// next token
func (p *Parser) next() lexer.Tokval {
	tok, ok := <-p.tokens
	if !ok {
		return tokEOF
	}
	return tok
}

// scry foretell the future using a crystal ball. Amount is how much
// of the future you want to foresee.
func (p *Parser) scry(amount int) []lexer.Tokval {
	if len(p.lookahead)+amount > 2 {
		panic("lookahead > 2")
	}

	sz := len(p.lookahead)
	for i := 0; i < amount-sz; i++ {
		val := p.next()
		p.lookahead = append(p.lookahead, val)
		if val.Type == token.EOF {
			break
		}
	}

	return p.lookahead
}

// forget what you had foresee
func (p *Parser) forget(amount int) {
	for i := 0; i < amount; i++ {
		p.lookahead = p.lookahead[1:]
	}
}

func parseIllegal(p *Parser) (ast.Node, error) {
	tok := p.lookahead[0]
	return nil, p.errorf(tok, "invalid token: %s",
		tok.Value)
}

func parseDecimal(p *Parser) (ast.Node, error) {
	tok := p.lookahead[0]
	defer p.forget(1)

	decstr := tok.Value
	f, err := strconv.ParseFloat(decstr.String(), 64)
	if err != nil {
		return nil, p.errorf(tok, err.Error())
	}
	return ast.NewNumber(f), nil
}

func parseHex(p *Parser) (ast.Node, error) {
	tok := p.lookahead[0]
	defer p.forget(1)

	hexstr := tok.Value
	prefix := utf16.S("0x")
	if hexstr.Index(prefix) == 0 {
		hexstr = hexstr.TrimPrefix(prefix)
	}
	hex, err := strconv.ParseInt(hexstr.String(), 16, 64)
	if err != nil {
		return nil, p.errorf(tok, err.Error())
	}
	return ast.NewIntNumber(hex), nil
}

// TODO(i4k): implement line and column of error
func (p *Parser) errorf(_ lexer.Tokval, f string, a ...interface{},
) error {
	return fmt.Errorf("%s:1:0: %s", p.filename, fmt.Sprintf(f, a...))
}

func isNumber(typ token.Type) bool {
	return typ == token.Decimal ||
		typ == token.Hexadecimal ||
		typ == token.Octal
}