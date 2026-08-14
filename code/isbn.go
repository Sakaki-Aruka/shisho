package code

import (
	"regexp"
	"strings"
)

func IsValidIsbn(code string) bool {
	return isValidIsbn10(code) || isValidIsbn13(code)
}

var replacer strings.Replacer = *strings.NewReplacer(" ", "-")

var isbn10Pattern regexp.Regexp = *regexp.MustCompile("([0-9]{1,5})[- ]?([0-9]+)[- ]?([0-9]+)[- ]?([0-9X])")
func isValidIsbn10(code string) bool {
	if !isbn10Pattern.MatchString(code) {
		return false
	}

	replaced := replacer.Replace(code)
	runes := []rune(replaced)
	sum := 0
	for i := 10; 2 <= i; i-- {
		sum += int(runes[i]) * i
	}

	calculatedCheckDigit := (11 - (sum % 11)) % 11
	checkDigitRune := runes[len(runes)-1]
	if calculatedCheckDigit < 10 {
		return calculatedCheckDigit == int(checkDigitRune)
	} else {
		return checkDigitRune == 'X'
	}
}

var isbn13Pattern regexp.Regexp = *regexp.MustCompile("97[89][- ]?([0-9]{1,5})[- ]?([0-9]+)[- ]?([0-9]+)[- ]?([0-9])")
func isValidIsbn13(code string) bool {
	if !isbn13Pattern.MatchString(code) {
		return false
	}

	replaced := replacer.Replace(code)
	runes := []rune(replaced)
	checkDigit := int(runes[len(runes)-1])
	sum := 0
	for i := 1; i <= 12; i++ {
		index := i - 1
		weight := 0
		if i % 2 == 0 {
			weight = 3
		} else {
			weight = 1
		}
		sum += int(runes[index]) * weight
	}
	return checkDigit == (10 - (sum % 10)) % 10
}