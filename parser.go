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
	}
	return output, err
}

func (p *Parser) ParseObject(obj map[string]interface{}) (interface{}, error) {
	var err error

	return obj, err
}

func (p *Parser) ParseArray(arr []interface{}) (interface{}, error) {
	var err error

	return arr, err
}
