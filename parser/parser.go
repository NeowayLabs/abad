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
var tokEOF = lexer.EOF

var (
	literalParsers   map[token.Type]parserfn
	unaryParsers     map[token.Type]parserfn
	varAssignParsers map[token.Type]parserfn
	nodeParsers      map[token.Type]parserfn
)

func init() {
	unaryParsers = map[token.Type]parserfn{
		token.Minus: parseUnary,
		token.Plus:  parseUnary,
	}
	literalParsers = map[token.Type]parserfn{
		token.Decimal:     parseDecimal,
		token.Hexadecimal: parseHex,
		token.String:      parseString,
		token.Bool:        parseBool,
		token.Undefined:   parseUndefined,
		token.Null:        parseNull,
	}
	varAssignParsers = mergeParsers(
		literalParsers,
		map[token.Type]parserfn{
			token.Ident: parseIdentExpr,
		},
	)
	nodeParsers = mergeParsers(
		literalParsers,
		unaryParsers,
		map[token.Type]parserfn{
			token.Ident: parseIdentExpr,
			token.Var:   parseVarDecls,
		},
	)
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

	// FIXME: This will probably not be enough to handle semicolon on the future
	for tok.Type == token.SemiColon {
		p.forget(1)
		p.scry(1)
		tok = p.lookahead[0]
	}

	if tok.Type == token.EOF {
		return nil, true, nil
	}

	if tok.Type == token.Illegal {
		_, err := parseIllegal(p)
		return nil, false, err
	}

	parser, ok := nodeParsers[tok.Type]

	if !ok {
		return nil, false, p.errorf(tok, "invalid token: %s", tok)
	}

	node, err := parser(p)
	if err != nil {
		return nil, false, err
	}

	// parsers should not leave tokens not processed
	// in the lookahead buffer.
	if len(p.lookahead) != 0 {
		return nil, false, fmt.Errorf(
			"parser for token[%v] not handled lookahead correctly, lookahead has[%v] but should be empty",
			tok,
			p.lookahead)
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
//
// Returns false if it reaches EOF before reading the desired amount
func (p *Parser) scry(amount int) bool {
	if len(p.lookahead)+amount > 2 {
		panic("lookahead > 2")
	}

	got := 0
	for i := 0; i < amount; i++ {
		val := p.next()
		p.lookahead = append(p.lookahead, val)
		got += 1
		if val.Type == token.EOF {
			break
		}
	}

	return got == amount
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

func parseString(p *Parser) (ast.Node, error) {
	tok := p.lookahead[0]
	defer p.forget(1)

	return ast.NewString(tok.Value), nil
}

func parseBool(p *Parser) (ast.Node, error) {
	tok := p.lookahead[0]
	defer p.forget(1)

	b, err := strconv.ParseBool(tok.Value.String())
	return ast.NewBool(b), err
}

func parseUndefined(p *Parser) (ast.Node, error) {
	p.forget(1)
	return ast.NewUndefined(), nil
}

func parseNull(p *Parser) (ast.Node, error) {
	p.forget(1)
	return ast.NewNull(), nil
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

func parseVarDecls(p *Parser) (ast.Node, error) {
	p.forget(1)
	return _parseVarDecls(p)
}

func _parseVarDecls(p *Parser) (ast.VarDecls, error) {

	if !p.scry(2) {
		return nil, fmt.Errorf("parser var decl: expected at least two tokens got[%s]", p.lookahead)
	}

	identifier := p.lookahead[0]
	possibleAssignment := p.lookahead[1]
	p.forget(2)

	if identifier.Type != token.Ident {
		return nil, fmt.Errorf("parser: var decl: expected identifier got[%s]", identifier)
	}

	varname := ast.NewIdent(identifier.Value)
	if possibleAssignment.Type == token.SemiColon {
		return ast.NewVarDecls(ast.NewVarDecl(varname, ast.NewUndefined())), nil
	}

	if possibleAssignment.Type != token.Assign {
		return nil, fmt.Errorf("parser: var decl: expected assignment token [=] got [%s]", possibleAssignment)
	}

	p.scry(1)
	assignExpr := p.lookahead[0]
	parser, hasparser := varAssignParsers[assignExpr.Type]

	if !hasparser {
		return nil, fmt.Errorf("parser: var decl: invalid token[%s] expected assigment expression", assignExpr)
	}

	val, err := parser(p)
	if err != nil {
		return nil, fmt.Errorf("parser: var decl: error[%s] parsing variable assign expression", err)
	}

	res := ast.NewVarDecls(ast.NewVarDecl(varname, val))
	p.scry(1)
	possibleSemiColon := p.lookahead[0]
	p.forget(1)

	if possibleSemiColon.Type == token.SemiColon || possibleSemiColon.Type == token.EOF {
		return res, nil
	}

	if possibleSemiColon.Type != token.Comma {
		return nil, fmt.Errorf("parser: var decl: invalid token[%s] expected comma", possibleSemiColon)
	}

	vars, err := _parseVarDecls(p)
	return append(res, vars...), err
}

func parseIdentExpr(p *Parser) (ast.Node, error) {
	tok := p.lookahead[0]
	p.scry(1)
	next := p.lookahead[1]

	// eg.: console.
	if next.Type == token.Dot {
		p.forget(1)
		return parseMemberExpr(p, ast.NewIdent(tok.Value))
	}

	// eg.: console(
	if next.Type == token.LParen {
		return parseCallExpr(p)
	}

	if next.Type != token.EOF && next.Type != token.SemiColon {
		return nil, p.errorf(next, "parser:identifier:unexpected token [%s]", next)
	}

	p.forget(2)

	return ast.NewIdent(tok.Value), nil
}

// state:
// lookahead[0] = token.Dot
func parseMemberExpr(p *Parser, object ast.Node) (ast.Node, error) {
	p.forget(1)

	tok := p.next()
	if tok.Type != token.Ident {
		return nil, p.errorf(tok, "unexpected %s", tok.Value)
	}

	member := ast.NewMemberExpr(object, ast.NewIdent(tok.Value))

	// TODO(i4k): Discuss this!
	// We can just return the MemberExpr here but then if next token is '('
	// then the main parser loop will need to retrieve the last parsed node
	// to use it as the object in the CallExpr.
	// Going into the CallExpr parser from here avoid to keep some states
	// in the parser structure but can lead to some code duplicates also.
	// I'm not sure of what's the best approach yet... Coding the naive one
	// for now.
	p.scry(1)

	tok = p.lookahead[0]
	if tok.Type == token.LParen {
		return parseMemberFuncall(p, member)
	}

	if tok.Type == token.Dot {
		return parseMemberExpr(p, member)
	}

	if tok.Type != token.EOF {
		return nil, p.errorf(tok, "unexpected %s", tok.Value)
	}

	p.forget(1)

	return member, nil
}

// state:
// lookahead[0] = token.LParen
func parseMemberFuncall(p *Parser, member *ast.MemberExpr) (ast.Node, error) {
	p.forget(1) // drops (
	args, err := parseFuncallArgs(p)
	if err != nil {
		return nil, err
	}

	return ast.NewCallExpr(member, args), nil
}

func parseFuncallArgs(p *Parser) ([]ast.Node, error) {
	if len(p.lookahead) != 0 {
		panic(fmt.Sprintf("parser: funcall args: unexpected non empty lookahead:%s", p.lookahead))
	}

	nextToken := func() lexer.Tokval {
		p.scry(1)
		return p.lookahead[0]
	}

	var args []ast.Node

	for {
		tok := nextToken()

		if tok.Type == token.RParen {
			p.forget(1)
			break
		}
		// TODO: not handling errors like successive commas
		if tok.Type == token.Comma {
			p.forget(1)
			continue
		}

		parser, hasParser := literalParsers[tok.Type]
		if hasParser {
			parsed, err := parser(p)
			if err != nil {
				return nil, err
			}
			args = append(args, parsed)
		} else {
			return nil, p.errorf(tok, "parser: funcall args: unexpected token [%s]", tok.Value)
		}
	}

	return args, nil
}

// state:
// lookahead[0] = token.Ident
// lookahead[1] = token.LParen
func parseCallExpr(p *Parser) (ast.Node, error) {
	ident := p.lookahead[0]
	p.forget(2) // drops <ident>(
	args, err := parseFuncallArgs(p)
	if err != nil {
		return nil, err
	}
	return ast.NewCallExpr(ast.NewIdent(ident.Value), args), nil
}

// TODO(i4k): implement line and column of error
func (p *Parser) errorf(_ lexer.Tokval, f string, a ...interface{}) error {
	return fmt.Errorf("%s:1:0: %s", p.filename, fmt.Sprintf(f, a...))
}

func mergeParsers(parsers ...map[token.Type]parserfn) map[token.Type]parserfn {
	res := map[token.Type]parserfn{}

	for _, parsermap := range parsers {
		for tokenType, fn := range parsermap {
			res[tokenType] = fn
		}
	}

	return res
}
