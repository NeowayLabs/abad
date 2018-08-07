package ast

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/token"
)

type (
	// Node type
	NodeType int

	// All node types implement the Node interface
	Node interface {
		Type() NodeType
		String() string
		Equal(other Node) bool
	}

	// Program Abstract Syntax Tree
	Program struct {
		Nodes []Node
	}

	Number float64

	String utf16.Str

	Bool bool

	Undefined struct{}

	Null struct{}

	// UnaryExpr is a unary expression (-a, +a, ~a, and so on)
	UnaryExpr struct {
		Operator token.Type
		Operand  Node
	}

	// MemberExpr handles get of object's properties
	// eg.: <object>.<property>
	MemberExpr struct {
		Object   Node
		Property Ident
	}

	CallExpr struct {
		Callee Node
		Args   []Node
	}

	// FunDecl is the syntatic function declaration
	FunDecl struct {
		Name Ident
		Args []Ident
		Body *Program
	}

	Ident utf16.Str

	VarDecl struct {
		Name  Ident
		Value Node
	}

	VarDecls []VarDecl
)

const (
	NodeProgram NodeType = iota + 1
	NodeFunDecl
	NodeVarDecl
	NodeVarDecls

	exprBegin

	NodeNumber
	NodeString
	NodeNull
	NodeUndefined
	NodeBool
	NodeUnaryExpr
	NodeMemberExpr
	NodeCallExpr
	NodeIdent

	exprEnd

	endNodeTypes
)

var nodeTypesNames = [...]string{
	NodeProgram:    "PROGRAM",
	NodeFunDecl:    "FUNDECL",
	NodeVarDecl:    "VARDECL",
	NodeVarDecls:   "VARDECLS",
	NodeNumber:     "NUMBER",
	NodeString:     "STRING",
	NodeBool:       "BOOLEAN",
	NodeUndefined:  "UNDEFINED",
	NodeNull:       "NULL",
	NodeUnaryExpr:  "UNARYEXPR",
	NodeMemberExpr: "MEMBEREXPR",
	NodeCallExpr:   "CALLEXPR",
	NodeIdent:      "IDENT",
	exprEnd:        "",
}

// console.log(Number.EPSILON);
// 2.220446049250313e-16
var ε = math.Pow(2, -52)

func (t NodeType) String() string {
	if t >= endNodeTypes {
		panic(fmt.Sprintf("unexpected node type: %d", t))
	}
	return nodeTypesNames[t]
}

func IsExpr(node Node) bool {
	return node.Type() > exprBegin &&
		node.Type() < exprEnd
}

func (_ *Program) Type() NodeType {
	return NodeProgram
}

func (p *Program) String() string {
	var stmts []string
	for _, stmt := range p.Nodes {
		stmts = append(stmts, stmt.String())
	}
	return strings.Join(stmts, "\n")
}

func (p *Program) Equal(other Node) bool {
	if other.Type() != NodeProgram {
		return false
	}

	o := other.(*Program)
	if len(p.Nodes) != len(o.Nodes) {
		return false
	}

	for i := 0; i < len(p.Nodes); i++ {
		if !p.Nodes[i].Equal(o.Nodes[i]) {
			return false
		}
	}
	return true
}

func NewString(a utf16.Str) String {
	return String(a)
}

func (String) Type() NodeType {
	return NodeString
}

func (s String) Equal(other Node) bool {
	otherStr, ok := other.(String)
	if !ok {
		return false
	}
	// Not efficient...but lazy
	return s.String() == otherStr.String()
}

func (s String) String() string {
	return utf16.Str(s).String()
}

func NewBool(a bool) Bool {
	return Bool(a)
}

func (a Bool) Equal(other Node) bool {
	otherb, ok := other.(Bool)
	if !ok {
		return false
	}
	return a == otherb
}

func (a Bool) String() string {
	return fmt.Sprintf("%t", a)
}

func (Bool) Type() NodeType {
	return NodeBool
}

func (b Bool) Value() bool {
	return bool(b)
}

func NewUndefined() Undefined {
	return Undefined{}
}

func (Undefined) Type() NodeType {
	return NodeUndefined
}

func (Undefined) Equal(other Node) bool {
	_, ok := other.(Undefined)
	return ok
}

func (Undefined) String() string {
	return "undefined"
}

func NewNull() Null {
	return Null{}
}

func (Null) Equal(other Node) bool {
	_, ok := other.(Null)
	return ok
}

func (Null) Type() NodeType {
	return NodeNull
}

func (Null) String() string {
	return "null"
}

func NewNumber(a float64) Number {
	return Number(a)
}

func NewIntNumber(a int64) Number {
	return Number(float64(a))
}

func (a Number) Value() float64 { return float64(a) }

func (_ Number) Type() NodeType {
	return NodeNumber
}

// TODO(i4k): Implements correct javascript number to string
// representation.
func (a Number) String() string {
	return strconv.FormatFloat(float64(a), 'f', -1, 64)
}

func (a Number) Equal(other Node) bool {
	if other.Type() != NodeNumber {
		return false
	}

	o := other.(Number)
	return floatEquals(float64(a), float64(o))
}

func NewUnaryExpr(operator token.Type, operand Node) *UnaryExpr {
	return &UnaryExpr{
		Operator: operator,
		Operand:  operand,
	}
}

func (_ *UnaryExpr) Type() NodeType {
	return NodeUnaryExpr
}

func (a *UnaryExpr) String() string {
	return fmt.Sprintf("%s%s", a.Operator, a.Operand)
}

func (a *UnaryExpr) Equal(other Node) bool {
	if other.Type() != a.Type() {
		return false
	}

	o := other.(*UnaryExpr)
	if a.Operator != o.Operator {
		return false
	}

	return a.Operand.Equal(o.Operand)
}

func NewIdent(ident utf16.Str) Ident {
	return Ident(ident)
}

func (_ Ident) Type() NodeType {
	return NodeIdent
}

func (an Ident) String() string {
	return utf16.Decode(utf16.Str(an))
}

func (an Ident) Equal(other Node) bool {
	if an.Type() != other.Type() {
		return false
	}

	astr := utf16.Str(an)
	ostr := utf16.Str(other.(Ident))

	if len(astr) != len(ostr) {
		return false
	}

	for i := 0; i < len(astr); i++ {
		if astr[i] != ostr[i] {
			return false
		}
	}

	return true
}

func NewMemberExpr(object Node, property Ident) *MemberExpr {
	return &MemberExpr{
		Object:   object,
		Property: property,
	}
}

func (m *MemberExpr) Type() NodeType { return NodeMemberExpr }
func (m *MemberExpr) String() string {
	return fmt.Sprintf("%s.%s", m.Object, m.Property)
}

func (m *MemberExpr) Equal(other Node) bool {
	if m.Type() != other.Type() {
		return false
	}

	o := other.(*MemberExpr)
	return m.Object.Equal(o.Object) &&
		m.Property.Equal(o.Property)
}

func NewVarDecl(name Ident, val Node) VarDecl {
	return VarDecl{
		Name:  name,
		Value: val,
	}
}

func (v VarDecl) Type() NodeType { return NodeVarDecl }

func (v VarDecl) Equal(other Node) bool {
	if other.Type() != v.Type() {
		return false
	}

	o := other.(VarDecl)
	return v.Name.Equal(o.Name) && v.Value.Equal(o.Value)
}

func (v VarDecl) String() string {
	return fmt.Sprintf("var %s = %s", v.Name, v.Value)
}

func NewVarDecls(vars ...VarDecl) VarDecls {
	return VarDecls(vars)
}

func (v VarDecls) Type() NodeType { return NodeVarDecls }

func (v VarDecls) Equal(other Node) bool {
	if other.Type() != v.Type() {
		return false
	}

	o := other.(VarDecls)

	if len(v) != len(o) {
		return false
	}

	for i, v := range v {
		if !v.Equal(o[i]) {
			return false
		}
	}

	return true
}

func (v VarDecls) String() string {
	varstr := []string{}
	for _, vardecl := range v {
		varstr = append(varstr, fmt.Sprintf("%s = %s", vardecl.Name, vardecl.Value))
	}
	return "var " + strings.Join(varstr, ",")
}

func NewCallExpr(callee Node, args []Node) *CallExpr {
	return &CallExpr{
		Callee: callee,
		Args:   args,
	}
}

func (c *CallExpr) Type() NodeType { return NodeCallExpr }
func (c *CallExpr) String() string {
	return fmt.Sprintf("%s(<args>)", c.Callee)
}

func (c *CallExpr) Equal(other Node) bool {
	if other.Type() != c.Type() {
		return false
	}

	o := other.(*CallExpr)

	if len(c.Args) != len(o.Args) {
		return false
	}

	for i := 0; i < len(c.Args); i++ {
		if !c.Args[i].Equal(o.Args[i]) {
			return false
		}
	}

	return c.Callee.Equal(o.Callee)
}

// NewFunDecl creates a new function declaration node.
func NewFunDecl(name Ident, args []Ident, body *Program) *FunDecl {
	return &FunDecl{
		Name: name,
		Args: args,
		Body: body,
	}
}

func (a *FunDecl) Type() NodeType {
	return NodeFunDecl
}

func (a *FunDecl) String() string {
	var args []string

	for _, arg := range a.Args {
		args = append(args, arg.String())
	}

	// TODO(i4k): improve identation
	return fmt.Sprintf("function %s(%s) {\n%s\n}",
		a.Name,
		strings.Join(args, ", "),
		a.Body.String(),
	)
}

func (a *FunDecl) Equal(other Node) bool {
	if other.Type() != NodeFunDecl {
		return false
	}

	o := other.(*FunDecl)

	if len(a.Args) != len(o.Args) {
		return false
	}

	for i := 0; i < len(a.Args); i++ {
		if !a.Args[i].Equal(o.Args[i]) {
			return false
		}
	}

	return a.Name.Equal(o.Name) && a.Body.Equal(o.Body)
}

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) < ε && math.Abs(b-a) < ε
}
