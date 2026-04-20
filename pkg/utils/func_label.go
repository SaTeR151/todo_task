package utils

import "fmt"

func AddFuncLabel(funcLabel string, err error) {
	if err != nil {
		err = fmt.Errorf("%s %v", funcLabel, err)
	}
}
