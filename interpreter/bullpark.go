package bullpark

import (
	"fmt"
	"strconv"
	"unicode"
)

type token struct {
	tokenType  string
	tokenValue string
}

type lexer struct {
	code        string
	currentChar byte
	index       int
	isEnd       bool
	reserve     map[string]token
}

// asasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasasas

type node struct {
	binOp    *binOp
	num      *num
	unaryOp  *unaryOp
	compound *compound
	assign   *assign
	_var     *_var
}

type binOp struct {
	left  *node
	op    token
	right *node
}

type unaryOp struct {
	token token
	expr  *node
}
type num struct {
	token token
}
type compound struct {
	children []*node
}
type assign struct {
	left  *node
	op    token
	right *node
}
type _var struct {
	token token
}

type parser struct {
	lexer        lexer
	currentToken token
}

func Parser(lexer *lexer) *parser {
	return &parser{lexer: *lexer}
}

func (par *parser) Init() {
	par.currentToken = par.lexer.get_next_token()
}

func (par *parser) eat(tokenType string) {
	// fmt.Println(par.currentToken.tokenType, ",", tokenType)
	if par.currentToken.tokenType == tokenType {
		par.currentToken = par.lexer.get_next_token()
		// fmt.Println(par.currentToken)
	} else {
		reportError("!!Token Type Error!!")
	}
}

func (par *parser) factor() *node {
	storeToken := par.currentToken
	if storeToken.tokenType == "PLUS" {
		par.eat("PLUS")
		return &node{binOp: nil, num: nil, unaryOp: &unaryOp{token: storeToken, expr: par.factor()}}
	} else if storeToken.tokenType == "MINUS" {
		par.eat("MINUS")
		return &node{binOp: nil, num: nil, unaryOp: &unaryOp{token: storeToken, expr: par.factor()}}
	} else if storeToken.tokenType == "INT" {
		par.eat("INT")
		return &node{binOp: nil, num: &num{token: storeToken}, unaryOp: nil}
	} else if storeToken.tokenType == "LEFTPAREN" {
		par.eat("LEFTPAREN")
		node := par.Expression()
		par.eat("RIGHTPAREN")
		return node
	} else {
		node := par.variable()
		return node
	}
}

func (par *parser) term() *node {
	_node := par.factor()
	var resultNode *node = _node
	for par.currentToken.tokenType == "MUL" || par.currentToken.tokenType == "DIV" {
		token := par.currentToken
		if token.tokenValue == "*" {
			par.eat("MUL")
		} else if token.tokenValue == "/" {
			par.eat("DIV")
		}
		resultNode = &node{binOp: &binOp{left: _node, op: token, right: par.factor()}, num: nil, unaryOp: nil}
	}
	return resultNode
}

func (par *parser) Expression() *node {
	_node := par.term()
	var resultNode *node = _node
	for par.currentToken.tokenType == "PLUS" || par.currentToken.tokenType == "MINUS" {
		token := par.currentToken
		if token.tokenValue == "+" {
			par.eat("PLUS")
		} else if token.tokenValue == "-" {
			par.eat("MINUS")
		}
		resultNode = &node{binOp: &binOp{left: _node, op: token, right: par.term()}, num: nil, unaryOp: nil}
	}
	return resultNode
}

func (par *parser) parse() *node {
	return par.program()
}

// ################################################################################################

// type nodeVisitor struct {
// }

// func (nv *nodeVisitor) visit(node any) {
// 	name := "visit_" + reflect.TypeOf(node).Name()
// 	visitor :=
// }

// ################################################################################################

type interpreter struct {
	parser       parser
	GLOBAL_SCOPE map[string]int
}

func Interpreter(parser *parser) *interpreter {
	return &interpreter{parser: *parser, GLOBAL_SCOPE: map[string]int{}}
}

func (inter *interpreter) visit_BinOp(node *node) int {
	if node.binOp.op.tokenType == "PLUS" {
		return inter.visit(node.binOp.left) + inter.visit(node.binOp.right)
	} else if node.binOp.op.tokenType == "MINUS" {
		return inter.visit(node.binOp.left) - inter.visit(node.binOp.right)
	} else if node.binOp.op.tokenType == "MUL" {
		return inter.visit(node.binOp.left) * inter.visit(node.binOp.right)
	} else if node.binOp.op.tokenType == "DIV" {
		return inter.visit(node.binOp.left) / inter.visit(node.binOp.right)
	}
	return 0
}

func (inter *interpreter) visit_Num(node *node) int {
	value, _ := strconv.Atoi(node.num.token.tokenValue)
	return value
}

func (inter *interpreter) visit_UnaryOp(node *node) int {
	op := node.unaryOp.token.tokenType
	if op == "PLUS" {
		return inter.visit(node.unaryOp.expr)
	} else if op == "MINUS" {
		return inter.visit(node.unaryOp.expr) * -1
	}
	return 0
}

func (inter *interpreter) visit_Compound(node *compound) {
	for i := 0; i < len(node.children); i++ {
		inter.visit(node.children[i])
	}
}

func (inter *interpreter) visit_Assign(node *node) {
	var_name := node.assign.left._var.token.tokenValue
	inter.GLOBAL_SCOPE[var_name] = inter.visit(node.assign.right)

	fmt.Println(var_name, inter.GLOBAL_SCOPE[var_name])
}

func (inter *interpreter) visit_Var(node *node) int {
	var_name := node._var.token.tokenValue
	val := inter.GLOBAL_SCOPE[var_name]
	return val
}

func (inter *interpreter) visit(node *node) int {
	if node.binOp != nil {
		return inter.visit_BinOp(node)
	} else if node.num != nil {
		return inter.visit_Num(node)
	} else if node.unaryOp != nil {
		return inter.visit_UnaryOp(node)
	} else if node.compound != nil {
		inter.visit_Compound(node.compound)
	} else if node.assign != nil {
		inter.visit_Assign(node)
	} else if node._var != nil {
		return inter.visit_Var(node)
	}
	return 0
}

func (inter *interpreter) Interpret() int {
	tree := inter.parser.parse()
	return inter.visit(tree)
}

// ################################################################################################

func reportError(content string) {
	panic(content)
}

func Lexer(code string) *lexer {
	return &lexer{code: code, index: 0, isEnd: false}
}

func (lex *lexer) Init() {
	lex.currentChar = lex.code[lex.index]
	lex.reserve = map[string]token{
		"var": {tokenType: "VARIABLE", tokenValue: "var"},
	}
}

func (lex *lexer) advance() {
	// CHECKING ABOUT NEXT CODE CHAR EXIST
	if (len(lex.code) > lex.index+1) && !lex.isEnd {
		lex.index += 1
		lex.currentChar = lex.code[lex.index]
	} else {
		lex.isEnd = true
	}
}

//pppppppppppppeekk
// func (lex *lexer) peek() byte {
// 	peekIndex := lex.index + 1
// 	if len(lex.code) > lex.index {
// 		lex.advance()
// 		return lex.code[peekIndex]
// 	}
// 	return ' '
// }

func (lex *lexer) id() token {
	var result string
	var reserveToken token
	for !lex.isEnd && unicode.IsLetter(rune(lex.currentChar)) {
		result += string(lex.currentChar)
		lex.advance()
	}
	val, exists := lex.reserve[result]
	if exists {
		reserveToken = token{tokenType: val.tokenType, tokenValue: val.tokenValue}
	} else {
		lex.reserve[result] = token{tokenType: "ID", tokenValue: result}
		reserveToken = token{tokenType: "ID", tokenValue: result}
	}
	return reserveToken
}

// func (lex *lexer) visitAssign(varName string) {
// 	lex.GLOBAL_SCOPE[varName] = lex.visit(node.right)
// }

func (lex *lexer) getInteger(char byte) string {
	_, err := strconv.Atoi(string(lex.currentChar))
	var result string
	for err == nil && (!lex.isEnd) {
		result += string(lex.currentChar)
		lex.advance()
		_, err = strconv.Atoi(string(lex.currentChar))
	}
	return result
}

func (lex *lexer) get_next_token() token {
	for !lex.isEnd {
		currentChar := lex.code[lex.index]
		switch currentChar {
		case ' ':
			lex.advance()
			continue
		case '+':
			lex.advance()
			return token{tokenType: "PLUS", tokenValue: string(currentChar)}
		case '-':
			lex.advance()
			return token{tokenType: "MINUS", tokenValue: string(currentChar)}
		case '*':
			lex.advance()
			return token{tokenType: "MUL", tokenValue: string(currentChar)}
		case '/':
			lex.advance()
			return token{tokenType: "DIV", tokenValue: string(currentChar)}
		case '(':
			lex.advance()
			return token{tokenType: "LEFTPAREN", tokenValue: string(currentChar)}
		case ')':
			lex.advance()
			return token{tokenType: "RIGHTPAREN", tokenValue: string(currentChar)}
		case '=':
			lex.advance()
			return token{tokenType: "ASSIGN", tokenValue: "="}
		case ';':
			lex.advance()
			return token{tokenType: "SEMI", tokenValue: ";"}
		case '.':
			lex.advance()
			return token{tokenType: "DOT", tokenValue: "."}
		default:
			_, err := strconv.Atoi(string(lex.currentChar))
			if err == nil { // literal (int)
				return token{tokenType: "INT", tokenValue: lex.getInteger(currentChar)}
			} else {
				return lex.id()
			}
		}
	}
	return token{"ERROR", ""}
}

// type interpreter struct {
// 	lexer        lexer
// 	currentToken token
// }

// func Interpreter(lexer *lexer) *interpreter {
// 	return &interpreter{lexer: *lexer}
// }

// func (inter *interpreter) Init() {
// 	inter.currentToken = inter.lexer.get_next_token()
// }

// func (inter *interpreter) eat(tokenType string) {
// 	if inter.currentToken.tokenType == tokenType {
// 		inter.currentToken = inter.lexer.get_next_token()
// 	} else {
// 		reportError("!!Token Type Error!!")
// 	}
// }

// func (inter *interpreter) factor() int {
// 	storeToken := inter.currentToken
// 	value, _ := strconv.Atoi(storeToken.tokenValue)
// 	if storeToken.tokenType == "INT" {
// 		inter.eat("INT")
// 		return value
// 	} else if inter.currentToken.tokenType == "LEFTPAREN" {
// 		inter.eat("LEFTPAREN")
// 		result := inter.Expression()
// 		inter.eat("RIGHTPAREN")
// 		return result
// 	}
// 	if storeToken.tokenType == "MINUS" {
// 		inter.eat("MINUS")
// 		result := inter.factor() * -1
// 		return result
// 	} else if storeToken.tokenType == "PLUS" {
// 		inter.eat("PLUS")
// 		result := inter.factor() * 1
// 		return result
// 	}
// 	return 6947
// }

// func (inter *interpreter) term() int {
// 	var result = inter.factor()
// 	for inter.currentToken.tokenType == "MUL" || inter.currentToken.tokenType == "DIV" {
// 		token := inter.currentToken

// 		if token.tokenValue == "*" {
// 			inter.eat("MUL")
// 			result *= inter.factor()
// 		} else if token.tokenValue == "/" {
// 			inter.eat("DIV")
// 			result /= inter.factor()
// 		}
// 	}
// 	return result
// }

// func (inter *interpreter) Expression() int {
// 	var result = inter.term()
// 	for inter.currentToken.tokenType == "PLUS" || inter.currentToken.tokenType == "MINUS" {
// 		token := inter.currentToken
// 		if token.tokenValue == "+" {
// 			inter.eat("PLUS")
// 			result += inter.term()
// 		} else if token.tokenValue == "-" {
// 			inter.eat("MINUS")
// 			result -= inter.term()
// 		}
// 	}
// 	return result
// }

func (par *parser) program() *node {
	result := par.compound_statement()
	par.eat("ERROR")
	return result
}

func (par *parser) compound_statement() *node {
	nodes := par.statement_list()

	root := &compound{children: []*node{}}
	for i := 0; i < len(nodes); i++ {
		root.children = append(root.children, nodes[i])
	}
	return &node{compound: root}
}

func (par *parser) statement_list() []*node {
	results := []*node{}
	results = append(results, par.statement())
	for par.currentToken.tokenType == "SEMI" {
		par.eat("SEMI")
		results = append(results, par.statement())
	}
	return results
}

func (par *parser) statement() *node {
	if par.currentToken.tokenType == "VARIABLE" {
		return par.assignment_statement()
	} else {
		return &node{}
	}
}

func (par *parser) assignment_statement() *node {
	par.eat("VARIABLE")
	left := par.variable()
	op := par.currentToken
	par.eat("ASSIGN")
	right := par.Expression()
	return &node{assign: &assign{left: left, op: op, right: right}}
}

func (par *parser) variable() *node {
	node := &node{_var: &_var{par.currentToken}}
	par.eat("ID")
	return node
}
