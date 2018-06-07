package itrml

import (
	"container/list"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var skip_ws_re = regexp.MustCompile("^\\s*")
var comma_re = regexp.MustCompile("^\\s*,")
var open_bracket_re = regexp.MustCompile("^\\s*\\[")
var close_bracket_re = regexp.MustCompile("^\\s*\\]")

var block_comment_re = regexp.MustCompile(
	"^\\s*\\[\\*" +

		"([^*]|\\*[^\\]])*" +

		"\\*\\]")

var first_quote_re = regexp.MustCompile("^\\s*\"")

var quoted_string_re = regexp.MustCompile("^\"((?:\\\\\\\\|\\\\\"|[^\"])*)\"")

var integer_re = regexp.MustCompile("^(?:-?\\d+)")
var long_re = regexp.MustCompile("^(?:-?\\d+)[lL]")
var double_re = regexp.MustCompile("^(?:-?\\d+\\.\\d+[dD]?|-?\\d+[dD])")
var word_re = regexp.MustCompile("^(?:[\\w\\.\\-:|]+)")

type Parser struct {
	Expression string
	Position   int
}

func first_quote(expr string, start int) []int {
	return first_quote_re.FindStringIndex(expr[start:])
}

func open_bracket(expr string, start int) []int {
	return open_bracket_re.FindStringIndex(expr[start:])
}

func close_bracket(expr string, start int) []int {
	return close_bracket_re.FindStringIndex(expr[start:])
}

func block_comment(expr string, start int) []int {
	return block_comment_re.FindStringIndex(expr[start:])
}

func (p *Parser) skip_ws() {
	loc := skip_ws_re.FindStringIndex(p.Expression[p.Position:])
	if loc != nil {
		p.Position += loc[1] - loc[0]
	}
}

func (p *Parser) skip_comma() {
	loc := comma_re.FindStringIndex(p.Expression[p.Position:])
	if loc != nil {
		p.Position += loc[1] - loc[0]
	}
}

func (p *Parser) parse() (interface{}, error) {
	p.skip_ws()

	for i := block_comment(p.Expression, p.Position); i != nil; i = block_comment(p.Expression, p.Position) {
		p.Position += i[1] - i[0]
	}

	i := open_bracket(p.Expression, p.Position)

	if i != nil {
		p.Position += i[1] - i[0]
		return p.parse_list()
	}

	if p.Position >= len(p.Expression) {
		return nil, nil
	}

	return p.parse_literal()
}

func (p *Parser) parse_list() (interface{}, error) {
	l := list.New()

	var i []int = nil

	// While there's no close bracket.
	for i = close_bracket(p.Expression, p.Position); i == nil; i = close_bracket(p.Expression, p.Position) {
		if v, err := p.parse(); err == nil {
			l.PushBack(v)
		} else {
			return nil, err
		}
		p.skip_comma()
		if p.Position >= len(p.Expression) {
			return nil, errors.New("Unclosed expression.")
		}
	}

	if i != nil {
		p.Position += i[1] - i[0]
	}

	return l, nil
}

func (p *Parser) parse_literal() (interface{}, error) {
	i := first_quote(p.Expression, p.Position)
	expr := p.Expression[p.Position:]

	if i != nil {
		if loc := quoted_string_re.FindStringSubmatchIndex(expr); loc != nil {
			p.Position += loc[1] - loc[0]
			tok := expr[loc[2]:loc[3]]
			tok = regexp.MustCompile("\\\\(.)").ReplaceAllString(tok, "$1")
			return tok, nil
		} else {
			return nil, errors.New("Unmatched \" starting at position " + string(p.Position))
		}
	}

	if loc := double_re.FindStringIndex(expr); loc != nil {
		tok := expr[loc[0]:loc[1]]
		p.Position += loc[1] - loc[0]
		if strings.HasSuffix(tok, "D") || strings.HasSuffix(tok, "d") {
			tok = tok[:len(tok)-1]
		}

		if f, err := strconv.ParseFloat(tok, 64); err == nil {
			return f, nil
		} else {
			return nil, err
		}
	}

	if loc := long_re.FindStringIndex(expr); loc != nil {
		tok := expr[loc[0]:loc[1]]
		p.Position += loc[1] - loc[0]
		if strings.HasSuffix(tok, "L") || strings.HasSuffix(tok, "l") {
			tok = tok[:len(tok)-1]
		}

		if l, err := strconv.ParseInt(tok, 10, 64); err == nil {
			return l, nil
		} else {
			return nil, err
		}
	}

	if loc := integer_re.FindStringIndex(expr); loc != nil {
		tok := expr[loc[0]:loc[1]]
		p.Position += loc[1] - loc[0]
		if n, err := strconv.ParseInt(tok, 10, 64); err == nil {
			return n, nil
		} else {
			return nil, err
		}
	}

	if loc := word_re.FindStringIndex(expr); loc != nil {
		tok := expr[loc[0]:loc[1]]
		p.Position += loc[1] - loc[0]
		return tok, nil
	}

	return nil, errors.New("Unexpected token at position " + string(p.Position))
}

func ParseExpression(e string) (interface{}, error) {
	p := Parser{e, 0}
	return p.parse()
}
