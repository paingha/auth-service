package utils

import (
	"strings"
)

//UppercaseName convert user first name and last name to have first character capitalized
func UppercaseName(name string) string {
	var emptyString []string
	newString := append(emptyString, strings.ToUpper(name[:1]), name[1:])
	//Join the array of string to make a single string
	concatedString := strings.Join(newString, "")
	return concatedString
}
