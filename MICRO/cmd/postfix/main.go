package postfix

import (
	"fmt"
	"strings"
	"unicode"
)

func findSequenceIndex(slice []string, sequence []string) int {
  for i := 0; i <= len(slice)-len(sequence); i++ {
    match := true
    for j := 0; j < len(sequence); j++ {
      if slice[i+j] != sequence[j] {
        match = false
        break
      }
    }
    if match {
      return i
    }
  }
  return -1
}

func ReplaceFirstSequence(slice []string, sequence []string, replacement string) []string {
  index := findSequenceIndex(slice, sequence)
  if index == -1 {
    return slice
  }
  newSlice := append(slice[:index], append([]string{replacement}, slice[index+len(sequence):]...)...)
  return newSlice
}

func isOperator(s string) bool {
  return s == "+" || s == "-" || s == "*" || s == "/"
}

func ToPostfix(exp string) ([]string, error) {
  var output []string
  var stack []string
  precedence := map[string]int{
    "+": 1,
    "-": 1,
    "*": 2,
    "/": 2,
  }

  newexp := exp
	for i := 0; i < len(exp) - 1; i++ {
		a := string(exp[i])
		b := string(exp[i + 1])
		if a == "-" && b != " " {
		j := i + 1
		for ; j < len(exp); j++ {
			if string(exp[j]) == " " || isOperator(string(exp[j])) {
			  break
			}
		}
    c := "(0 - " + exp[i + 1:j] + ")"
		newexp = strings.Replace(newexp, a + exp[i + 1:j], c, 1)
		}
	}
	exp = newexp

  greaterPrecedence := func(op1, op2 string) bool {
    return precedence[op1] >= precedence[op2]
  }
  i := 0
  for i < len(exp) {
    ch := rune(exp[i])

    if unicode.IsSpace(ch) {
      // Ignore spaces
      i++
      continue
    }

    if unicode.IsDigit(ch) || ch == '.' {
      start := i
      for i < len(exp) && (unicode.IsDigit(rune(exp[i])) || exp[i] == '.') {
        i++
      }
      output = append(output, exp[start:i])
      continue
    }

    if ch == '(' {
      stack = append(stack, string(ch))
      i++
      continue
    }

    if ch == ')' {
      for len(stack) > 0 && stack[len(stack)-1] != "(" {
        output = append(output, stack[len(stack)-1])
        stack = stack[:len(stack)-1]
      }
      if len(stack) == 0 || stack[len(stack)-1] != "(" {
        return nil, fmt.Errorf("mismatched parentheses")
      }
      stack = stack[:len(stack)-1]
      i++
      continue
    }
    if isOperator(string(ch)) {
      for len(stack) > 0 && isOperator(stack[len(stack)-1]) && greaterPrecedence(stack[len(stack)-1], string(ch)) {
        output = append(output, stack[len(stack)-1])
        stack = stack[:len(stack)-1]
      }
      stack = append(stack, string(ch))
      i++
      continue
    }

    return nil, fmt.Errorf("invalid character: %c", ch)
  }
  for len(stack) > 0 {
    if stack[len(stack)-1] == "(" || stack[len(stack)-1] == ")" {
      return nil, fmt.Errorf("mismatched parentheses")
    }
    output = append(output, stack[len(stack)-1])
    stack = stack[:len(stack)-1]
  }
  return output, nil
}