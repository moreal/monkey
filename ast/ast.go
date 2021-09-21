package ast

import (
	"bytes"
	"fmt"
	"github.com/moreal/monkey/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	fmt.Stringer
	statementNode()
}

type Expression interface {
	Node
	fmt.Stringer
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (*LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	out.WriteString(ls.Value.String())
	out.WriteString(token.SEMICOLON)

	return out.String()
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (*ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	out.WriteString(rs.Value.String())
	out.WriteString(token.SEMICOLON)

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (*Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (*ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	var out bytes.Buffer

	out.WriteString(es.String())
	out.WriteString(token.SEMICOLON)

	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (*IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (*PrefixExpression) expressionNode() {}
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}
func (p *PrefixExpression) String() string {
	return fmt.Sprintf("%s%s%s%s", token.LPAREN, p.Operator, p.Right.String(), token.RPAREN)
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (*InfixExpression) expressionNode() {}
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}
func (i *InfixExpression) String() string {
	return fmt.Sprintf("%s%s %s %s%s", token.LPAREN, i.Left.String(), i.Operator, i.Right.String(), token.RPAREN)
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (*Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (*IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (*BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (*FunctionLiteral) expressionNode() {}
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	var params []string
	for _, param := range f.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(f.TokenLiteral())
	out.WriteString(token.LPAREN)
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(token.RPAREN + " ")
	out.WriteString(f.Body.String())

	return out.String()
}
