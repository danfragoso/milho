package interpreter

import (
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
	var acc int64
	for _, n := range numbers {
		acc -= n
	}

	return &Result{
		Type:  Number,
		Value: strconv.FormatInt(acc, 10),
	}, nil
}
