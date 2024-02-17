package Services

import (
	"216/internal/orchestrator/Database"
	"216/internal/orchestrator/Entities"
	"fmt"
	"github.com/gofrs/uuid"
	"log"
	"regexp"
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

	for _, char := range expression {
		if char == '(' {
			countOpenBrackets++
		} else if char == ')' {
			countClosedBrackets++
			if countClosedBrackets > countOpenBrackets {
				return false
			}
		} else if !unicode.IsDigit(char) && !IsOperator(string(char)) {
			// Если символ не является цифрой и не является одним из допустимых операторов, выражение некорректно
			return false
		}
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

func ToPostfix(expression string) string {
	outputQueue := []string{}
	operatorStack := []string{}

	tokens := regexp.MustCompile(`([\d.]+|[\+\-\*\/\(\)])`).FindAllString(expression, -1)

	for _, token := range tokens {
		if token == "" {
			continue
		} else if token == "(" {
			operatorStack = append(operatorStack, token)
		} else if token == ")" {
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != "(" {
				outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			operatorStack = operatorStack[:len(operatorStack)-1]
		} else if IsOperator(token) {
			for len(operatorStack) > 0 && Priority(operatorStack[len(operatorStack)-1]) >= Priority(token) {
				outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			operatorStack = append(operatorStack, token)
		} else {
			outputQueue = append(outputQueue, token)
		}
	}
	for len(operatorStack) > 0 {
		outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
		operatorStack = operatorStack[:len(operatorStack)-1]
	}
	return strings.Join(outputQueue, " ")
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
					Previous:  &b,
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
					Previous:  &b,
					Next:      &a,
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
func searchUuid(id uuid.UUID, arr []Task) bool {
	for _, v := range arr {
		if v.Id == id {
			return true
			break
		}
	}
	return false
}

// ////
func IsValidArithmeticExpression(expression string) bool {
	// Удаляем все пробелы из выражения
	expression = strings.ReplaceAll(expression, " ", "")

	// Проверяем, что выражение не пустое
	if len(expression) == 0 {
		return false
	}

	// Проверяем первый символ выражения
	firstChar := rune(expression[0])
	if !unicode.IsDigit(firstChar) && firstChar != '(' {
		return false
	}

	// Проверяем последний символ выражения
	lastChar := rune(expression[len(expression)-1])
	if !unicode.IsDigit(lastChar) && lastChar != ')' {
		return false
	}

	// Проверяем парность скобок
	var stack []rune
	for _, ch := range expression {
		if ch == '(' {
			stack = append(stack, ch)
		} else if ch == ')' {
			if len(stack) == 0 {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	if len(stack) > 0 {
		return false
	}

	// Проверяем последовательность символов
	for i := 0; i < len(expression)-1; i++ {
		currChar := rune(expression[i])
		nextChar := rune(expression[i+1])

		if unicode.IsDigit(currChar) {
			if !unicode.IsDigit(nextChar) && nextChar != ')' {
				return false
			}
		} else if currChar == '(' {
			if !unicode.IsDigit(nextChar) && nextChar != '(' {
				return false
			}
		} else if currChar == ')' {
			if !isOperator(nextChar) && nextChar != ')' {
				return false
			}
		} else if isOperator(currChar) {
			if !unicode.IsDigit(nextChar) && nextChar != '(' {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func isOperator(ch rune) bool {
	operators := map[rune]bool{
		'+': true,
		'-': true,
		'*': true,
		'/': true,
	}
	return operators[ch]
}
