package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"regexp"
)

import "github.com/golang-collections/collections/stack" 

type RuleType int

const (
	CLASS_COMMENT_RULE 		RuleType = 1
	VARIABLE_NAME_RULE 		RuleType = 2
	METHOD_COMMENT_RULE 	RuleType = 13
	METHOD_LENGTH_RULE 		RuleType = 23
	VARIABLE_COMMENT_RULE	RuleType = 33
	INDENT_RULE 			RuleType = 4
	MAGIC_NUMBER_RULE 		RuleType = 5
	CONSTANT_RULE 			RuleType = 6
	LINE_LENGTH_RULE 		RuleType = 7
	CVS_RULE 				RuleType = 8
	IMPORT_RULE 			RuleType = 9
	EXCEPTION_RULE 			RuleType = 10
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Il manque le fichier. Attention, c'est une déduction de 3 points sur la note finale.")
		os.Exit(1)
	}

	fileName := os.Args[1]

	data, err := ioutil.ReadFile(fileName)
    if err != nil {
        fmt.Println("Problème lors de la lecture du fichier.", err)
		os.Exit(1)
	}

	array_str := strings.Split(string(data), "\n")

	isInMethod := false
	methodName := ""
	methodLength := 0
	stack := stack.New()

	reMethodName := regexp.MustCompile(`\s([a-zA-Z]+\s?\(.*\))`)
	reVariableName := regexp.MustCompile(`\w+\s([^ijk<>!]|\w{2,})(?:\s?=\s?(?:.)+)\;`)
	//reMagicNumber := regexp.MustCompile(`\w+\s([^ijk<>!]|\w{2,})(?:\s?=\s?(?:[2-9])+)\;`)


	lineBefore := ""

	for lineNb, line := range array_str {
		lineNb += 1

		if strings.Contains(line, "class") {
			if !(strings.Contains(line, "*") || strings.Contains(line, "//")) {
				if !(strings.Contains(lineBefore, "*") || strings.Contains(lineBefore, "//")) {
					smite(METHOD_COMMENT_RULE, lineNb, line)
				}
			}
		}

		/*
		if reMagicNumber.MatchString(line) {
			smite(MAGIC_NUMBER_RULE, lineNb, line)
		}
		*/

		if reVariableName.MatchString(line) {
			if !strings.Contains(line, "//") {
				smite(VARIABLE_COMMENT_RULE, lineNb, line)
			}
		}

		if len(line) >= 150 {
			smite(LINE_LENGTH_RULE, lineNb, line)
		}

		if strings.Contains(line, "\t") {
			smite(INDENT_RULE, lineNb, line)
		}

		if !isInMethod {
			match, _ := regexp.MatchString(`\)\s?\{`, line)
			if match {
				methodName = reMethodName.FindStringSubmatch(line)[1]
				isInMethod = true
				stack.Push("{")
			}
		} else {
			methodLength++

			for _, char := range line {

				if char == '}' {
					stack.Pop()
				} else if char == '{' {
					stack.Push("{")
				}

				if stack.Len() == 0 {
					isInMethod = false
					
					if methodLength >= 50 {
						smite2(METHOD_LENGTH_RULE, lineNb, line, methodName)
					}

					methodLength = 0
					methodName = ""
				}
			}
		}

		lineBefore = line;


	}
}

func smite(rule RuleType, lineNb int, lineStr string) {
	smite2(rule, lineNb, lineStr, "")
}

func smite2(rule RuleType, lineNb int, lineStr string, msg string) {
	ruleDescription := ""

	switch rule {
		case METHOD_LENGTH_RULE:
			ruleDescription = "The method is longer than 50 lines."
		case LINE_LENGTH_RULE:
			ruleDescription = "Line is longer than 120 characters."
		case INDENT_RULE:
			ruleDescription = "There is a tab at this line."
		case METHOD_COMMENT_RULE:
			ruleDescription = "The class or variable doesn't have any comment."
		case VARIABLE_COMMENT_RULE:
			ruleDescription = "The variable is not commented"
		case MAGIC_NUMBER_RULE:
			ruleDescription = "This value isn't a constant, it may be a magic number."
	}

	fmt.Println("COMMANDMENT", rule % 10, ": " + ruleDescription)
	if (rule == METHOD_LENGTH_RULE) {
		fmt.Println("In Method '" + msg + "'")
	}
	
	fmt.Println("At line", lineNb)
	fmt.Println("·", lineStr)
	fmt.Println("You will be punished.")
	fmt.Println("=======================================")


}