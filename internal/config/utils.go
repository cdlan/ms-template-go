package config

import (
	"log"
	"net/url"
	"strconv"
)

func stringToUrl(str string) url.URL {

	Url, err := url.Parse(str)
	if err != nil {
		log.Println(err)
	}

	return *Url
}

func stringToBool(str string) bool {

	boolVal, err := strconv.ParseBool(str)
	if err != nil {
		log.Println(err.Error())
	}

	return boolVal
}

func stringToInt(str string) int {

	intVal, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err.Error())
	}

	return intVal
}
