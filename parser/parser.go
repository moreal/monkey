package parser

import (
	"github.com/moreal/monkey/ast"
	"github.com/moreal/monkey/lexer"
	"github.com/moreal/monkey/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.LTE:      LESSGREATER,
	token.GT:       LESSGREATER,
	token.GTE:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	parser := &Parser{l: l, prefixParseFns: make(map[token.TokenType]prefixParseFn), infixParseFns: make(map[token.TokenType]infixParseFn)}

	parser.registerPrefixParseFn(token.IDENT, parser.parseIdentifier)
	parser.registerPrefixParseFn(token.INT, parser.parseIntgerLiteral)

	parser.registerPrefixParseFn(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefixParseFn(token.BANG, parser.parsePrefixExpression)

	parser.registerInfixParseFn(token.PLUS, parser.parseInfixExpression)
	parser.registerInfixParseFn(token.MINUS, parser.parseInfixExpression)
	parser.registerInfixParseFn(token.SLASH, parser.parseInfixExpression)
	parser.registerInfixParseFn(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfixParseFn(token.LT, parser.parseInfixExpression)
	parser.registerInfixParseFn(token.LTE, parser.parseInfixExpression)
	parser.registerInfixParseFn(token.GT, parser.parseInfixExpression)
	parser.registerInfixParseFn(token.GTE, parser.parseInfixExpression)
	parser.registerInfixParseFn(token.EQ, parser.parseInfixExpression)
	parser.registerInfixParseFn(token.NEQ, parser.parseInfixExpression)

	parser.nextToken()
	parser.nextToken()
	return parser
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	letStmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: Parse expression

	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return letStmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	// TODO: Parse expression

	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	// TODO: Parse expression

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}

	leftExp := prefix()

	for p.peekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntgerLiteral() ast.Expression {
	integerLiteral := &ast.IntegerLiteral{Token: p.curToken}

	v, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		return nil
	}

	integerLiteral.Value = v
	return integerLiteral
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekToken.Type == tokenType {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) registerPrefixParseFn(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixParseFn(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
