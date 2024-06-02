package main

import (
	"bufio"
	"fmt"
	"io"
	bullpark "main/go/interpreter"
	"os"
	"strings"
)

func scanString(s *bufio.Scanner) (string, error) {
	if s.Scan() {
		return s.Text(), nil
	}
	err := s.Err()
	if err == nil {
		err = io.EOF
	}
	return "", err
}

func main() {
	for {
		s := bufio.NewScanner(os.Stdin)

		var code string
		for {
			fmt.Print("inter> ")
			var err error
			code, err = scanString(s)
			if err == nil {
				code := strings.TrimSpace(code)
				min := 1
				if len(code) >= min {
					break
				}
			}
			fmt.Println("input error:", err)
		}
		lexer := bullpark.Lexer(code)
		lexer.Init()
		parser := bullpark.Parser(lexer)
		parser.Init()
		interpreter := bullpark.Interpreter(parser)
		result := interpreter.Interpret()
		fmt.Println(result)
	}
}
