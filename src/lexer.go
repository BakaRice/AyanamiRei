package AyanamiRei

import (
	"fmt"
	"regexp"
	"strings"
)

//token const
const (
	TOKEN_EOF         = iota // end-of-file
	TOKEN_VER_PREFIX         // $
	TOKEN_LEFT_PAREN         //(
	TOKEN_RIGHT_PAREN        //)
	TOKEN_EQUAL              //=
	TOKEN_QUOTE              //"
	TOKEN_DUQUOTE            //""
	TOKEN_NAME               //Name ::= [_A-Za-z][_0-9A-Za-z]*
	TOKEN_PRINT              //print
	TOKEN_IGNORED            //Ignored
)

var tokenNameMap = map[int]string{
	TOKEN_EOF:         "EOF",   // end-of-file
	TOKEN_VER_PREFIX:  "$",     // $
	TOKEN_LEFT_PAREN:  "(",     //(
	TOKEN_RIGHT_PAREN: ")",     //)
	TOKEN_EQUAL:       "=",     //=
	TOKEN_QUOTE:       "\"",    //"
	TOKEN_DUQUOTE:     "\"\"",  //""
	TOKEN_NAME:        "Name",  //Name ::= [_A-Za-z][_0-9A-Za-z]*
	TOKEN_PRINT:       "print", //print
	TOKEN_IGNORED:     "Ignored",
}

var keywords = map[string]int{
	"print": TOKEN_PRINT,
}

var regexName = regexp.MustCompile(`^[_\d\w]+`)

type Lexer struct {
	sourceCode       string //源代码 直接读取源代码文件并输入进来
	lineNum          int    //记录当前执行到的代码的行数
	nextToken        string //下一个token
	nextTokenType    int    //下一个token的类型
	nextTokenLineNum int    //下一个token的行数
}

/*
本质上Lexer是一个状态机，他只要能处理当前状态和跳到下一个状态，就可以一直工作下去
*/
func newLexer(sourceCode string) *Lexer {
	return &Lexer{sourceCode, 1, "", 0, 0}
}

/**
用于断言下一个Token是什么，并且由于内部执行了GetNextToken(),游标会自动向前移动
*/
func (lexer *Lexer) NextTokenIs(tokenType int) (lineNum int, token string) {
	nowLineNum, nowTokenType, nowToken := lexer.GetNextToken()
	//syntax error
	if tokenType != nowTokenType {
		err := fmt.Sprintf("NextTokenIs():syntax error near '%s',expected token:{%s} but got {%s}",
			tokenNameMap[nowTokenType], tokenNameMap[tokenType], tokenNameMap[nowTokenType])
		panic(err)
	}
	return nowLineNum, nowToken
}

/**
相对于lookAhead(),多了一个参数，仍然是先看一下下一个Token是什么,如果跟输入的相同就会跳过，如果不同则不会跳过，
这个函数是为了Ignored Token特别定制的
*/
func (lexer *Lexer) LookAheadAndSkip(expectedType int) {
	//get next token
	nowLineNum := lexer.lineNum
	lineNum, tokenType, token := lexer.GetNextToken()
	//not is expected type,reverse cursor
	if tokenType != expectedType {
		lexer.lineNum = nowLineNum
		lexer.nextTokenLineNum = lineNum
		lexer.nextTokenType = tokenType
		lexer.nextToken = token
	}
}

/**
用于返回下一个token是什么，不过它并不会将游标向前移动（准确说是移动了，又移了回来）
*/
func (lexer *Lexer) LookAhead() int {
	//lexer.nextToken * already set
	if lexer.nextTokenLineNum > 0 {
		return lexer.nextTokenType
	}
	//set it
	nowLineNum := lexer.lineNum
	lineNum, tokenType, token := lexer.GetNextToken()
	lexer.lineNum = nowLineNum
	lexer.nextTokenLineNum = lineNum
	lexer.nextTokenType = tokenType
	lexer.nextToken = token
	return tokenType
}

func (lexer *Lexer) nextSourceCodeIs(s string) bool {
	return strings.HasPrefix(lexer.sourceCode, s)
}

func (lexer *Lexer) skipSourceCode(n int) {
	lexer.sourceCode = lexer.sourceCode[n:]
}

func (lexer *Lexer) isIgnored() bool {
	isIgorned := false
	//target pattern
	isNewLine := func(c byte) bool {
		return c == '\r' || c == '\n'
	}
	isWhiteSpace := func(c byte) bool {
		switch c {
		case '\t', '\n', '\v', '\f', '\r', ' ':
			return true
		}
		return false
	}
	//match
	for len(lexer.sourceCode) > 0 {
		if lexer.nextSourceCodeIs("\r\n") || lexer.nextSourceCodeIs("\n\r") {
			lexer.skipSourceCode(2)
			lexer.lineNum += 1
			isIgorned = true
		} else if isNewLine(lexer.sourceCode[0]) {
			lexer.skipSourceCode(1)
			lexer.lineNum += 1
			isIgorned = true
		} else if isWhiteSpace(lexer.sourceCode[0]) {
			lexer.skipSourceCode(1)
			isIgorned = true
		} else {
			break
		}
	}
	return isIgorned
}

func (lexer *Lexer) scan(regexp *regexp.Regexp) string {
	if token := regexp.FindString(lexer.sourceCode); token != "" {
		lexer.skipSourceCode(len(token))
		return token
	}
	panic("unreachable!")
	return ""
}

func (lexer *Lexer) scanBeforeToken(token string) string {
	s := strings.Split(lexer.sourceCode, token)
	if len(s) < 2 {
		panic("unreachable!")
		return ""
	}
	lexer.skipSourceCode(len(s[0]))
	return s[0]
}

func (lexer *Lexer) scanName() string {
	return lexer.scan(regexName)
}

func (lexer *Lexer) GetNextToken() (lineNum int, tokenType int, token string) {
	//next token already loaded
	if lexer.nextTokenLineNum > 0 {
		lineNum = lexer.nextTokenLineNum
		tokenType = lexer.nextTokenType
		token = lexer.nextToken
		lexer.lineNum = lexer.nextTokenLineNum
		lexer.nextTokenLineNum = 0
		return
	}
	return lexer.MatchToken()
}

/**
实现查看当前字符是什么token ：查看当前字符，然后识别是什么token
*/
func (lexer *Lexer) MatchToken() (lineNum int, tokenType int, token string) {
	//check ignored
	if lexer.isIgnored() {
		return lexer.lineNum, TOKEN_IGNORED, "Ignored"
	}
	//finsh
	if len(lexer.sourceCode) == 0 {
		return lexer.lineNum, TOKEN_EOF, tokenNameMap[TOKEN_EOF]
	}
	//check token 实现取出当前字符，然后用case进行匹配token
	switch lexer.sourceCode[0] {
	case '$':
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_VER_PREFIX, "$"
	case '(':
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_LEFT_PAREN, "("
	case ')':
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_RIGHT_PAREN, ")"
	case '=':
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_EQUAL, "="
	case '"':
		if lexer.nextSourceCodeIs("\"\"") {
			return lexer.lineNum, TOKEN_DUQUOTE, "\"\""
		}
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_QUOTE, "\""
	}
	//check multiple character token
	//如果是下划线或者字母开头的就把name扫描出来
	if lexer.sourceCode[0] == '_' || isLetter(lexer.sourceCode[0]) {
		token := lexer.scanName()
		if tokenType, isMatch := keywords[token]; isMatch {
			return lexer.lineNum, tokenType, token
		} else {
			return lexer.lineNum, TOKEN_NAME, token
		}
	}
	//unexpected symbol
	err := fmt.Sprintf("MatchToken():unexpected symbol near '%q'.", lexer.sourceCode[0])
	panic(err)
	return
}

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}
