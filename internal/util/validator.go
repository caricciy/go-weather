package util

import "regexp"

func CheckCEPIsValid(cep string) bool {
	//there is no way to have a error in this function, so we can ignore the error
	check, _ := regexp.MatchString(`^[0-9]{8}$`, cep)
	return check
}
