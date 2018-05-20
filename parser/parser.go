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
	Value: utf16.S("EOF"),
}

var (
	keywordParsers = map[token.Type]parserfn{}
	literalParsers = map[token.Type]parserfn{
		token.Decimal:     parseDecimal,
		token.Hexadecimal: parseHex,
	}
	unaryParsers map[token.Type]parserfn
)

func init() {
	unaryParsers = map[token.Type]parserfn{
		token.Minus: parseUnary,
		token.Plus:  parseUnary,
	}
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
		node, eof, err := p.parseNode()
		if err != nil {
			return nil, err
		}

		if eof {
			break
		}

		nodes = append(nodes, node)
	}

	return &ast.Program{
		Nodes: nodes,
	}, nil
}

func (p *Parser) parseNode() (n ast.Node, eof bool, err error) {
	p.scry(1)

	tok := p.lookahead[0]
	if tok.Type == token.EOF {
		return nil, true, nil
	}

	if tok.Type == token.Illegal {
		_, err := parseIllegal(p)
		return nil, false, err
	}

	var parser parserfn
	var hasparser bool

	for _, parsers := range []map[token.Type]parserfn{
		keywordParsers,
		literalParsers,
		unaryParsers,
		map[token.Type]parserfn{
			token.Ident: parseIdentExpr,
		},
	} {
		parser, hasparser = parsers[tok.Type]
		if hasparser {
			break
		}
	}

	if !hasparser {
		return nil, false, p.errorf(tok, "invalid token: %s", tok)
	}

	node, err := parser(p)
	if err != nil {
		return nil, false, err
	}

	// parsers should not leave tokens not processed
	// in the lookahead buffer.
	if len(p.lookahead) != 0 {
		panic("parsers not handling lookahead correctly")
	}
	return node, false, nil
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

	for i := 0; i < amount; i++ {
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
	p.lookahead = p.lookahead[amount:]
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
	hexPrefix := utf16.S("0x")
	hexstr = hexstr.TrimPrefix(hexPrefix)
	hex, err := strconv.ParseInt(hexstr.String(), 16, 64)
	if err != nil {
		return nil, p.errorf(tok, err.Error())
	}

	return ast.NewIntNumber(hex), nil
}

func parseUnary(p *Parser) (ast.Node, error) {
	tok := p.lookahead[0]
	if !token.IsUnaryOperator(tok.Type) {
		return nil, p.errorf(tok, "unexpected: %s", tok.Type)
	}
	p.forget(1)
	expr, eof, err := p.parseNode()
	if err != nil {
		return nil, err
	}

	if eof {
		return nil, p.errorf(tok, "unexpected eof")
	}

	if !ast.IsExpr(expr) {
		return nil, p.errorf(tok, "expected expression, but got %s",
			expr.Type())
	}

	return ast.NewUnaryExpr(tok.Type, expr), nil
}

func parseIdentExpr(p *Parser) (ast.Node, error) {
	tok := p.lookahead[0]
	p.scry(1)
	next := p.lookahead[1]

	// eg.: console.
	if next.Type == token.Dot {
		return parseMemberExpr(p)
	}

	// eg.: console(
	if next.Type == token.LParen {
		return parseCallExpr(p)
	}

	if next.Type != token.EOF {
		return nil, p.errorf(next, "unexpected token '%s'", next)
	}

	p.forget(2)

	return ast.NewIdent(tok.Value), nil
}

// state:
// lookahead[0] = token.Ident
// lookahead[1] = token.Dot
func parseMemberExpr(p *Parser) (ast.Node, error) {
	object := ast.NewIdent(p.lookahead[0].Value)
	p.forget(2)

	tok := p.next()
	if tok.Type != token.Ident {
		return nil, p.errorf(tok, "unexpected %s", tok.Value)
	}

	return ast.NewMemberExpr(object, ast.NewIdent(tok.Value)), nil
}

// state:
// lookahead[0] = token.Ident
// lookahead[1] = token.LParen
func parseCallExpr(p *Parser) (ast.Node, error) {
	return nil, nil
}

// TODO(i4k): implement line and column of error
func (p *Parser) errorf(_ lexer.Tokval, f string, a ...interface{},
) error {
	return fmt.Errorf("%s:1:0: %s", p.filename, fmt.Sprintf(f, a...))
}
