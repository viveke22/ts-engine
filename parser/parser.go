package parser

import (
	"fmt"
	"strconv"
	"strings"
	"ts-engine/ast"
	"ts-engine/lexer"
	"ts-engine/token"
)

const (
	_ int = iota
	LOWEST
	ASSIGN      // =
	LOGICAL_OR  // ||
	LOGICAL_AND // &&
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.EQ:            EQUALS,
	token.NOT_EQ:        EQUALS,
	token.LT:            LESSGREATER,
	token.GT:            LESSGREATER,
	token.PLUS:          SUM,
	token.MINUS:         SUM,
	token.SLASH:         PRODUCT,
	token.ASTERISK:      PRODUCT,
	token.MOD:           PRODUCT,
	token.EQ_STRICT:     EQUALS,
	token.NOT_EQ_STRICT: EQUALS,
	token.AND:           LOGICAL_AND,
	token.OR:            LOGICAL_OR,
	token.LPAREN:        CALL,
	token.DOT:           CALL,
	token.LBRACKET:      INDEX,
	token.ASSIGN:        ASSIGN,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
	Strict         bool
}

func New(l *lexer.Lexer, strict bool) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
		Strict: strict,
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.AWAIT, p.parsePrefixExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.EQ_STRICT, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ_STRICT, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.DOT, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignmentExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET, token.CONST, token.VAR:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.EXPORT:
		return p.parseExportStatement()
	case token.DECLARE:
		return p.parseDeclareStatement()
	case token.IMPORT:
		return p.parseImportStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExportStatement() *ast.ExportStatement {
	stmt := &ast.ExportStatement{Token: p.curToken}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Optional type annotation: let x: number = ...
	if p.peekTokenIs(token.COLON) {
		p.nextToken() // consume COLON
		stmt.Name.Type = p.parseTypeAnnotation()
	} else if p.Strict {
		msg := fmt.Sprintf("missing type annotation for variable '%s' in strict mode (.ts file)", stmt.Name.Value)
		p.errors = append(p.errors, msg)
		return nil
	}

	// Declaration without assignment: let x: number;
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return stmt
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// Handle empty return (return;)
	if p.curTokenIs(token.SEMICOLON) {
		return stmt
	}

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
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

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	exp := &ast.AssignmentExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Value = p.parseExpression(LOWEST)

	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = p.curToken.Literal
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type: function(): void { ... }
	if p.peekTokenIs(token.COLON) {
		p.nextToken() // consume COLON
		lit.ReturnType = p.parseTypeAnnotation()
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Optional type annotation: (x: number)
	if p.peekTokenIs(token.COLON) {
		p.nextToken() // consume COLON
		ident.Type = p.parseTypeAnnotation()
	}

	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		// Optional type annotation: (..., y: string)
		if p.peekTokenIs(token.COLON) {
			p.nextToken() // consume COLON
			ident.Type = p.parseTypeAnnotation()
		}

		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IF) {
			p.nextToken()
			expression.Alternative = p.parseIfExpression()
		} else {
			if !p.expectPeek(token.LBRACE) {
				return nil
			}
			expression.Alternative = p.parseBlockStatement()
		}
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
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

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseImportStatement() ast.Statement {
	stmt := &ast.ImportStatement{Token: p.curToken}

	// Expect '*'
	if !p.expectPeek(token.ASTERISK) {
		return nil
	}

	// Expect 'as'
	if !p.expectPeek(token.AS) {
		return nil
	}

	// Expect Identifier (alias)
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Alias = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Expect 'from'
	if !p.expectPeek(token.FROM) {
		return nil
	}

	// Expect String (source)
	if !p.expectPeek(token.STRING) {
		return nil
	}

	stmt.Source = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

	// Optional Semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseDeclareStatement() ast.Statement {
	// declare var x: any;
	// Just consume until semicolon
	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}
	// Do NOT consume the semicolon here; ParseProgram loop does p.nextToken()
	// which expects curToken to be the last token of the statement (e.g. semicolon).

	// Return nil effectively ignores this statement in the AST
	return nil
}

func (p *Parser) parseTypeAnnotation() string {
	// If the next token is [, it might be a Tuple (or Array start if we supported [T] syntax, which we do now for tuples)
	if p.peekTokenIs(token.LBRACKET) {
		p.nextToken() // consume [
		return p.parseTupleType()
	}

	p.nextToken() // consume first type identifier

	typeName := p.curToken.Literal

	for p.peekTokenIs(token.DOT) {
		p.nextToken() // consume DOT
		p.nextToken() // consume next part of type
		typeName += "." + p.curToken.Literal
	}

	// Handle Array Types: number[] or string[][]
	for p.peekTokenIs(token.LBRACKET) {
		p.nextToken() // consume [
		if !p.peekTokenIs(token.RBRACKET) {
			break
		}
		p.nextToken() // consume ]
		typeName += "[]"
	}

	return typeName
}

func (p *Parser) parseTupleType() string {
	// curToken is [
	var types []string

	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
		return "[]"
	}

	for {
		// Parse the inner type
		// parseTypeAnnotation expects to be called *before* the type token
		// Since we are at [ or , (previous token), and peek is the start of type.
		// e.g. [ string
		// cur=[, peek=string.
		// parseTypeAnnotation calls nextToken -> cur=string. Correct.

		typeStr := p.parseTypeAnnotation()
		types = append(types, typeStr)

		if p.peekTokenIs(token.RBRACKET) {
			break
		}

		if !p.expectPeek(token.COMMA) {
			return ""
		}
	}

	if !p.expectPeek(token.RBRACKET) {
		return ""
	}

	return "[" + strings.Join(types, ", ") + "]"
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}
