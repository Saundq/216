package Services

import (
	"216/internal/orchestrator/Database"
	"216/internal/orchestrator/Entities"
	"fmt"
	"github.com/gofrs/uuid"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func IsOperator(op string) bool {
	return op == "+" || op == "-" || op == "*" || op == "/"
}

func Priority(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return 0
}

func IsValidExpression(expression string) bool {
	countOpenBrackets := 0
	countClosedBrackets := 0
	countWhiteSpaces := 0
	expression = strings.TrimSpace(expression)

	for _, char := range expression {
		if char == '(' {
			countOpenBrackets++
		} else if char == ')' {
			countClosedBrackets++
			if countClosedBrackets > countOpenBrackets {
				return false
			}
		} else if unicode.IsSpace(char) {
			countWhiteSpaces++
		} else if !unicode.IsDigit(char) && !IsOperator(string(char)) && !unicode.IsSpace(char) {
			return false
		}
	}
	if string(expression[0]) == "+" || string(expression[0]) == "*" || string(expression[0]) == "/" {
		return false
	}
	if countWhiteSpaces == 0 {
		return false
	}

	last := expression[len(expression)-1]
	if string(last) == "(" || IsOperator(string(last)) {
		return false
	}

	if countOpenBrackets != countClosedBrackets {
		return false
	}

	return true
}

func ToPostfix(infix string) string {
	var result strings.Builder
	stack := make([]rune, 0)
	tokens := strings.Fields(infix)

	for _, token := range tokens {
		if len(token) > 1 && token[0] == '-' {
			num, _ := strconv.Atoi(token)
			result.WriteString(fmt.Sprintf("%d ", num))
		} else if IsOperator(string(token[0])) {
			for len(stack) > 0 && Priority(string(token[0])) <= Priority(string(stack[len(stack)-1])) {
				result.WriteString(string(stack[len(stack)-1]) + " ")
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, rune(token[0]))
		} else if token == "(" {
			stack = append(stack, '(')
		} else if token == ")" {
			for stack[len(stack)-1] != '(' {
				result.WriteString(string(stack[len(stack)-1]) + " ")
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]
		} else {
			result.WriteString(token + " ")
		}
	}

	for len(stack) > 0 {
		result.WriteString(string(stack[len(stack)-1]) + " ")
		stack = stack[:len(stack)-1]
	}

	return strings.TrimSpace(result.String())
}

func EvaluatePostfix(postfix string) float64 {
	stack := []float64{}
	tokens := strings.Fields(postfix)
	var operations []Entities.ArithmeticOperation
	Database.Instance.Find(&operations)
	operationsMap := make(map[string]int)

	for _, v := range operations {
		operationsMap[v.Value] = v.LeadTime
	}

	for _, token := range tokens {
		if !IsOperator(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else {
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			switch token {
			case "+":
				time.Sleep(time.Duration(operationsMap[token]) * time.Second)
				stack = append(stack, a+b)
			case "-":
				time.Sleep(time.Duration(operationsMap[token]) * time.Second)
				stack = append(stack, a-b)
			case "*":
				time.Sleep(time.Duration(operationsMap[token]) * time.Second)
				stack = append(stack, a*b)
			case "/":
				time.Sleep(time.Duration(operationsMap[token]) * time.Second)
				stack = append(stack, a/b)
			}
		}
	}
	return stack[0]
}

type Task struct {
	Id        uuid.UUID
	Full      bool
	Number    bool
	Operand1  float64
	Operand2  float64
	Operation string
	Previous  *Task
	Next      *Task
	Polsk     string
}

func Mathslice(postfix string) []Task {
	//stack := []float64{}
	tokens := strings.Fields(postfix)
	result := []Task{}
	result1 := []uuid.UUID{}

	for _, token := range tokens {
		uid, _ := uuid.NewV4()
		if !IsOperator(token) {
			num, _ := strconv.ParseFloat(token, 64)
			//stack = append(stack, num)
			//uid, _ := uuid.NewV4()
			result = append(result, Task{
				Id:       uid,
				Full:     false,
				Operand1: num,
				Number:   true,
			})
			result1 = append(result1, uid)

		} else {
			log.Println(token)
			b := result[len(result)-1]
			a := result[len(result)-2]
			if a.Number && b.Number {
				//log.Println("minus tut", token)
				result = result[:len(result)-2]
				result = append(result, Task{
					Id:        uid,
					Operation: token,
					Full:      true,
					Number:    false,
					Operand1:  a.Operand1,
					Operand2:  b.Operand1,
					Polsk:     fmt.Sprintf("%.2f", a.Operand1) + " " + fmt.Sprintf("%.2f", b.Operand1) + " " + token,
				})
				result1 = result1[:len(result1)-2]
				result1 = append(result1, uid)
			} else if !a.Number && b.Number {
				result = result[:len(result)-1] //!!!!!
				result = append(result, Task{
					Id:        uid,
					Operation: token,
					Full:      false,
					Operand1:  b.Operand1,
					Polsk:     fmt.Sprintf("%.2f", b.Operand1) + " " + token,
					Previous:  &a,
					Number:    false,
				})
				result1 = result1[:len(result1)-2]
				result1 = append(result1, uid)
			} else if a.Number && !b.Number {
				index := len(result) - 2
				result = append(result[:index], result[index+1:]...)
				//result = result[:len(result)-1] //!!!!!
				result = append(result, Task{
					Id:        uid,
					Operation: token,
					Full:      false,
					Operand1:  a.Operand1,
					Polsk:     fmt.Sprintf("%.2f", a.Operand1) + " " + token,
					Next:      &b,
					Number:    false,
				})
				result1 = result1[:len(result1)-2]
				result1 = append(result1, uid)
			} else {
				//result = result[:len(result)-1]
				if result1[len(result1)-1] != b.Id {
					//if searchUuid(result1[len(result1)-1], result) {
					b.Id = result1[len(result1)-1]
					//}
				}
				if result1[len(result1)-2] != a.Id {
					//if searchUuid(result1[len(result1)-2], result) {
					a.Id = result1[len(result1)-2]
					//}
				}
				result = append(result, Task{
					Id:        uid,
					Operation: token,
					Full:      false,
					Polsk:     fmt.Sprintf("%.2f", a.Operand1) + " " + token,
					Previous:  &a,
					Next:      &b,
					Number:    false,
				})
				result1 = result1[:len(result1)-2]
				result1 = append(result1, uid)
			}

		}
		log.Println(result1)
		log.Println(result)
	}
	return result
}
