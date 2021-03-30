package interpreter

import (
	"errors"
	"fmt"
	"strconv"
)

func evalTypeNumber(n *Result) (int64, error) {
	if n.Type != Number {
		return 0, fmt.Errorf("expected parameter type Number, got %s", n.Type)
	}

	if n.Value == "" {
		return 0, nil
	}

	parsedInt, err := strconv.ParseInt(n.Value, 10, 64)
	if err != nil {
		return 0, err
	}

	return parsedInt, nil
}

func numberPrepareParams(params []*Result) ([]int64, error) {
	var nParams []int64
	for _, parameter := range params {
		n, err := evalTypeNumber(parameter)
		if err != nil {
			return nil, err
		}

		nParams = append(nParams, n)
	}

	return nParams, nil
}

func numberSum(numbers []int64) (*Result, error) {
	var acc int64
	for _, n := range numbers {
		acc += n
	}

	return &Result{
		Type:  Number,
		Value: strconv.FormatInt(acc, 10),
	}, nil
}

func numberSub(numbers []int64) (*Result, error) {
	if len(numbers) == 0 {
		return nil, errors.New("Wrong number of args '0' passed to Number:[-] function")
	}

	var acc int64
	if len(numbers) == 1 {
		acc = -numbers[0]
	} else {
		acc = numbers[0]
		for _, n := range numbers[1:] {
			acc -= n
		}
	}

	return &Result{
		Type:  Number,
		Value: strconv.FormatInt(acc, 10),
	}, nil
}

func numberMul(numbers []int64) (*Result, error) {
	var acc int64 = 1
	for _, n := range numbers {
		acc *= n
	}

	return &Result{
		Type:  Number,
		Value: strconv.FormatInt(acc, 10),
	}, nil
}

func numberDiv(numbers []int64) (*Result, error) {
	if len(numbers) == 0 {
		return nil, errors.New("Wrong number of args '0' passed to Number:[-] function")
	}

	var acc int64 = 1
	if len(numbers) > 1 {
		acc = numbers[0]
		for _, n := range numbers[1:] {
			if n == 0 {
				return nil, errors.New("Divide by zero error")
			}

			acc /= n
		}
	}

	return &Result{
		Type:  Number,
		Value: strconv.FormatInt(acc, 10),
	}, nil
}
