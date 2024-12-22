package calculator

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	Number TokenType = iota
	Operator
	LeftParen
	RightParen
)

type Token struct {
	Type  TokenType
	Value string
}

// с помощью польской аннотации переводим строку в постфикс
func tokenize(expression string) ([]Token, error) {
	var tokens []Token // здесь мапа tokens используется как стэк
	var number strings.Builder

	for _, char := range expression {
		if unicode.IsDigit(char) || char == '.' {
			number.WriteRune(char)
		} else {
			if number.Len() > 0 {
				tokens = append(tokens, Token{Type: Number, Value: number.String()})
				number.Reset()
			}
			switch char {
			case '+', '-', '*', '/':
				tokens = append(tokens, Token{Type: Operator, Value: string(char)})
			case '(':
				tokens = append(tokens, Token{Type: LeftParen, Value: string(char)})
			case ')':
				tokens = append(tokens, Token{Type: RightParen, Value: string(char)})
			case ' ':
				continue
			default:
				return nil, errors.New("неизвестный символ")
			}
		}
	}

	if number.Len() > 0 {
		tokens = append(tokens, Token{Type: Number, Value: number.String()})
	}

	return tokens, nil
}

// в этой функции мы расставляем приоритет операциям
// сложение и вычитание получают приоритет 1, когда умножение и деление получают большие приоритет
func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

// функция математических операций (сложение, вычитание, деление, умножение)
func eval(a, b float64, op string) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("деление на ноль")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("неизвестная операция %s", op)
	}
}

// здесь просиходит сам процесс решение выражения
// метод берет данные из стека и решает выражение, опираясь на приоритет знаков операции и расстоновку скобок
// обязательно: если в выражении есть скобки, то они идут парами! не должно быть незакрытых скобок
func calculate(tokens []Token) (float64, error) {
	var values []float64
	var ops []string

	for _, token := range tokens {
		switch token.Type {
		case Number:
			num, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				log.Fatal(err)
				return 0, err
			}
			values = append(values, num)

		case Operator:
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(token.Value) {
				if len(values) < 2 {
					return 0, errors.New("не хватает значений для операции")
				}
				b := values[len(values)-1]
				values = values[:len(values)-1]
				a := values[len(values)-1]
				values = values[:len(values)-1]
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]

				result, err := eval(a, b, op)
				if err != nil {
					if err.Error() == "деление на ноль" {
						return 0, fmt.Errorf("деление на ноль")
					}
					log.Fatal(err)
					return 0, err
				}
				values = append(values, result)
			}
			ops = append(ops, token.Value)

		case LeftParen:
			ops = append(ops, token.Value)

		case RightParen:
			for len(ops) > 0 && ops[len(ops)-1] != "(" {
				if len(values) < 2 {
					return 0, errors.New("не хватает значений для операции")
				}
				b := values[len(values)-1]
				values = values[:len(values)-1]
				a := values[len(values)-1]
				values = values[:len(values)-1]
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]

				result, err := eval(a, b, op)
				if err != nil {
					if err.Error() == "деление на ноль" {
						return 0, fmt.Errorf("деление на ноль")
					}
					log.Fatal(err)
					return 0, err
				}
				values = append(values, result)
			}
			if len(ops) == 0 || ops[len(ops)-1] != "(" {
				return 0, errors.New("несоответствующая скобка")
			}
			ops = ops[:len(ops)-1]

		default:
			return 0, errors.New("неизвестный оператор")
		}
	}

	for len(ops) > 0 {
		if len(values) < 2 {
			return 0, errors.New("не хватает значений для операции")
		}
		b := values[len(values)-1]
		values = values[:len(values)-1]
		a := values[len(values)-1]
		values = values[:len(values)-1]
		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]

		result, err := eval(a, b, op)
		if err != nil {
			if err.Error() == "деление на ноль" {
				return 0, fmt.Errorf("деление на ноль")
			}
			log.Fatal(err)
			return 0, err
		}
		values = append(values, result)
	}

	if len(values) != 1 {
		return 0, errors.New("ошибка в выражении")
	}
	return values[0], nil
}

// функция объединяет в себе функию перевода строки в постфикс и функцию калькулятора, чтобы по итогу сразу вывести ответ
// в виде числа float64
func Calc(expression string) (float64, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		if err.Error() == "неизвестный символ" {
			return 0, fmt.Errorf("неизвестный символ")
		}
		log.Fatal(err)
		return 0, err
	}
	return calculate(tokens)
}
