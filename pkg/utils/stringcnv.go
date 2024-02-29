package utils

import (
	"net/url"
	"strconv"
)

func StringToUrl(str string) (url.URL, error) {

	Url, err := url.Parse(str)
	if err != nil {
		return url.URL{}, err
	}

	return *Url, nil
}

func StringToBool(str string) (bool, error) {

	boolVal, err := strconv.ParseBool(str)
	if err != nil {
		return false, err
	}

	return boolVal, nil
}

func StringToInt(str string) (int, error) {

	intVal, err := strconv.Atoi(str)
	if err != nil {
		return -1, err
	}

	return intVal, nil
}