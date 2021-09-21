package ast

import (
	"fmt"
	"github.com/moreal/monkey/token"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
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

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (*ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
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
