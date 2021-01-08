package AyanamiRei

import "fmt"

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
)

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
实现查看当前字符是什么token ：查看当前字符，然后识别是什么token
*/
func (lexer *Lexer) MatchToken() (lineNum int, tokenType int, token string) {
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
		//if lexer.nextSourceIs("\"\"") {
		//	return lexer.lineNum, TOKEN_DUQUOTE, "\"\""
		//}
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_QUOTE, "\""
	}
	//unexpected symbol
	err := fmt.Sprintf("MatchToken():unexpected symbol near '%q'.", lexer.sourceCode[0])
	panic(err)
	return
}

func (lexer *Lexer) skipSourceCode(n int) {
	lexer.sourceCode = lexer.sourceCode[n:]
}

//func (lexer *Lexer) nextSourceIs(s string) bool {
//
//}
