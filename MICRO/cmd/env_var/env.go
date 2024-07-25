package envvar

import (
	"os"
	"strconv"
)

func CompPower() int {
	compPower := os.Getenv("COMPUTING_POWER")
	gocount, err := strconv.Atoi(compPower)
	if compPower == "" || err != nil {
		return 10
	}
	return gocount
}

func OperVars(op string) int {
	if op == "+" {
		addenv := os.Getenv("TIME_ADDITION_MS")
		res, err := strconv.Atoi(addenv)
		if err != nil || addenv == "" {
			return 100
		}
		return res
	} else if op == "-" {
		subenv := os.Getenv("TIME_SUBTRACTION_MS")
		res, err := strconv.Atoi(subenv)
		if err != nil || subenv == "" {
			return 100
		}
		return res
	} else if op == "*" {
		mulenv := os.Getenv("TIME_MULTIPLICATIONS_MS")
		res, err := strconv.Atoi(mulenv)
		if err != nil || mulenv == "" {
			return 100
		}
		return res
	} else {
		divenv := os.Getenv("TIME_DIVISIONS_MS")
		res, err := strconv.Atoi(divenv)
		if err != nil || divenv == "" {
			return 100
		}
		return res
	}
}