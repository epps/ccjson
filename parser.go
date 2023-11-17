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
	switch tok.Type {
	case True:
		output = true
	case False:
		output = false
	case Null:
		output = nil
	case String:
		output = tok.Literal
	case Number:
		output, err = strconv.ParseFloat(tok.Literal, 64)
	case BeginObject:
		output, err = p.ParseObject(make(map[string]interface{}))
	case BeginArray:
		output, err = p.ParseArray(make([]interface{}, 0))
	case Illegal:
		err = fmt.Errorf("illegal token %s encountered", tok.Literal)
	case EOF:
		err = fmt.Errorf("unexpected end of input")
	default:
		err = fmt.Errorf("unknown token %s", tok.Literal)
	}

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

	// To iterative parse multi-key objects, we loop until we don't find encounter
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
		var value interface{}
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
				return obj, err
			}
		case BeginObject:
			value, err = p.ParseObject(make(map[string]interface{}))
		case BeginArray:
			value, err = p.ParseArray(make([]interface{}, 0))
		default:
			return obj, fmt.Errorf("invalid value: %s", tok.Literal)
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

	if tok.Type != EndArray {
		err = fmt.Errorf("expected end of input but found %s", tok.Literal)
	}

	return arr, err
}
