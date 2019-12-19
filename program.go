package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/golang-collections/collections/stack"
)

type RuleType int

const (
	ClassCommentRule    RuleType = 1
	VariableNameRule    RuleType = 2
	MethodCommentRule   RuleType = 13
	MethodLengthRule    RuleType = 23
	VariableCommentRule RuleType = 33
	IndentRule          RuleType = 4
	MagicNumberRule     RuleType = 5
	ConstantRule        RuleType = 6
	LineLengthRule      RuleType = 7
	CvsRule             RuleType = 8
	ImportRule          RuleType = 9
	ExceptionRule       RuleType = 10
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Missing file. Warning, this substracts 3 points from the final grade.")
		os.Exit(1)
	}

	fileName := os.Args[1]

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading file.", err)
		os.Exit(1)
	}

	arrayStr := strings.Split(string(data), "\n")

	isInMethod := false
	methodName := ""
	methodLength := 0
	stack := stack.New()

	reMethodName := regexp.MustCompile(`\s([a-zA-Z]+\s?\(.*\))`)
	reVariableName := regexp.MustCompile(`\w+\s([^ijk<>!]|\w{2,})(?:\s?=\s?(?:.)+)\;`)
	// reMagicNumber := regexp.MustCompile(`\w+\s([^ijk<>!]|\w{2,})(?:\s?=\s?(?:[2-9])+)\;`)
	// test

	lineBefore := ""

	for lineNb, line := range arrayStr {
		lineNb++

		if strings.Contains(line, "class") {
			if !(strings.Contains(line, "*") || strings.Contains(line, "//")) {
				if !(strings.Contains(lineBefore, "*") || strings.Contains(lineBefore, "//")) {
					smite(MethodCommentRule, lineNb, line)
				}
			}
		}

		/*
			if reMagicNumber.MatchString(line) {
				smite(MagicNumberRule, lineNb, line)
			}
		*/

		if reVariableName.MatchString(line) {
			if !strings.Contains(line, "//") {
				smite(VariableCommentRule, lineNb, line)
			}
		}

		if len(line) >= 150 {
			smite(LineLengthRule, lineNb, line)
		}

		if strings.Contains(line, "\t") {
			smite(IndentRule, lineNb, line)
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
						smite2(MethodLengthRule, lineNb, line, methodName)
					}

					methodLength = 0
					methodName = ""
				}
			}
		}

		lineBefore = line
	}
}

func smite(rule RuleType, lineNb int, lineStr string) {
	smite2(rule, lineNb, lineStr, "")
}

func smite2(rule RuleType, lineNb int, lineStr string, msg string) {
	ruleDescription := ""

	switch rule {
	case MethodLengthRule:
		ruleDescription = "The method is longer than 50 lines."
	case LineLengthRule:
		ruleDescription = "The line is longer than 120 characters."
	case IndentRule:
		ruleDescription = "This line is indented with tabs."
	case MethodCommentRule:
		ruleDescription = "The class or variable doesn't have any comment."
	case VariableCommentRule:
		ruleDescription = "The variable is not commented"
	case MagicNumberRule:
		ruleDescription = "This value isn't a constant, it may be a magic number."
	}

	fmt.Println("COMMANDMENT", rule%10, ": ", ruleDescription)
	if rule == MethodLengthRule {
		fmt.Println("In Method '" + msg + "'")
	}

	fmt.Println("At line", lineNb)
	fmt.Println("·", lineStr)
	fmt.Println("You will be punished.")
	fmt.Println("=======================================\n")
}
