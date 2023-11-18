package main

import (
	"fmt"
	"strconv"
)

func NewParser(input string) *Parser {
	return &Parser{
		lexer: NewLexer(input),
	}
}

type Parser struct {
	lexer *Lexer
}

func (p *Parser) Parse() (interface{}, error) {
	var output interface{}
	var err error

	tok := p.lexer.NextToken()

	output, err = p.ParseToken(tok)

	tok = p.lexer.NextToken()

	if tok.Type != EOF {
		err = fmt.Errorf("expected end of input but found %s", tok.Literal)
	}
	return output, err
}

func (p *Parser) ParseObject(obj map[string]interface{}) (interface{}, error) {
	var err error
	// The BeginObject token is read in Parse, and given it has no use in our parsing
	// other than to identify an object value, we immediately advance to the next token
	tok := p.lexer.NextToken()

	if tok.Type == EndObject {
		return obj, err
	}

	// To iteratively parse multi-key objects, we loop until we don't find encounter
	// the ValueSeparator token
	for {
		if tok.Type != String {
			return obj, fmt.Errorf("expected key but found %s", tok.Literal)
		}

		key := tok.Literal

		tok = p.lexer.NextToken()

		if tok.Type != NameSeparator {
			return obj, fmt.Errorf("expected name separate but found %s", tok.Literal)
		}

		tok = p.lexer.NextToken()

		value, err := p.ParseToken(tok)
		if err != nil {
			return obj, err
		}

		obj[key] = value

		tok = p.lexer.NextToken()

		if tok.Type != ValueSeparator {
			break
		}

		// We advance past the ValueSeparator token to arrive at the next
		// key in the object
		tok = p.lexer.NextToken()
	}

	if tok.Type != EndObject {
		err = fmt.Errorf("expected } but found %s", tok.Literal)
	}

	return obj, err
}

func (p *Parser) ParseArray(arr []interface{}) (interface{}, error) {
	var err error

	tok := p.lexer.NextToken()

	if tok.Type == EndArray {
		return arr, err
	}

	// To iteratively parse arrays, we loop until we don't encounter
	// the ValueSeparator token
	for {
		value, err := p.ParseToken(tok)
		if err != nil {
			return arr, err
		}

		arr = append(arr, value)

		tok = p.lexer.NextToken()

		if tok.Type != ValueSeparator {
			break
		}

		// We advance past the ValueSeparator token to arrive at the next
		// value in the array
		tok = p.lexer.NextToken()
	}

	if tok.Type != EndArray {
		err = fmt.Errorf("expected end of input but found %s", tok.Literal)
	}

	return arr, err
}

func (p *Parser) ParseToken(tok Token) (interface{}, error) {
	var value interface{}
	var err error
	switch tok.Type {
	case True:
		value = true
	case False:
		value = false
	case Null:
		value = nil
	case String:
		value = tok.Literal
	case Number:
		value, err = strconv.ParseFloat(tok.Literal, 64)
		if err != nil {
			return value, err
		}
	case BeginObject:
		value, err = p.ParseObject(make(map[string]interface{}))
	case BeginArray:
		value, err = p.ParseArray(make([]interface{}, 0))
	case Illegal:
		err = fmt.Errorf("illegal token %s encountered", tok.Literal)
	case EOF:
		err = fmt.Errorf("unexpected end of input")
	default:
		err = fmt.Errorf("unknown token %s", tok.Literal)
	}
	return value, err
}
