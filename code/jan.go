package code

import (
	"errors"
	"regexp"
	"strconv"
)

var pattern regexp.Regexp = *regexp.MustCompile("^192([0-9]{4})([0-9]{5})([0-9])$")

var invalidCodeError error = errors.New("Invalid BookJan second line code")

func IsValidBookJan(code string) bool {
	if !pattern.MatchString(code) {
		return false
	} 

	runes := []rune(code)

	oddSum, evenSum := 0, 0
	
	for i := 0; i <= 11; i++ {
		num := int(runes[i])
		if i % 2 == 0 {
			oddSum += num
		} else {
			evenSum += num
		}
	}

	checkDigit := int(runes[12])
	return (10 - (((evenSum * 3) + oddSum) % 10)) % 10 == checkDigit
}

func GetJanClassificationCode(secondLine string) (int, error) {
	if !IsValidBookJan(secondLine) {
		return 0, invalidCodeError
	}
	
	matches := pattern.FindStringSubmatch(secondLine)
	if result, err := strconv.Atoi(matches[1]); err != nil {
		return 0, invalidCodeError
	} else {
		return result, nil
	}
}

func GetPrice(secondLine string) (int, error) {
	if !IsValidBookJan(secondLine) {
		return 0, invalidCodeError
	}

	matches := pattern.FindStringSubmatch(secondLine)
	if result, err := strconv.Atoi(matches[2]); err != nil {
		return 0, invalidCodeError
	} else {
		return result, nil
	}
}

