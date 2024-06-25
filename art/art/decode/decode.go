package decode

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrMalformedInput = errors.New("malformed input")
)

func DecodeArt(encodedString string) (string, error) {
	// Adjusted regex to match encoded sequences correctly.
	correctPattern := regexp.MustCompile(`\[(\d+)\s([^]]+)]`)
	encodingAttempt := regexp.MustCompile(`\[.*?]`)

	// Checking if the string contains unmatched or malformed encoding attempts
	if encodingAttempt.MatchString(encodedString) && !correctPattern.MatchString(encodedString) {
		return "", ErrMalformedInput
	}

	decodedString := correctPattern.ReplaceAllStringFunc(encodedString, func(match string) string {
		parts := correctPattern.FindStringSubmatch(match)
		count, err := strconv.Atoi(parts[1])
		if err != nil {
			// This should never happen because the regex ensures that a number is matched.
			return match
		}
		// Replicates the specified character or string 'count' times.
		return strings.Repeat(parts[2], count)
	})

	// If the encoding attempt is present but the string remains unchanged, consider it malformed.
	if encodingAttempt.MatchString(encodedString) && decodedString == encodedString {
		return "", ErrMalformedInput
	}

	return decodedString, nil
}

func DecodeMultiLine(encodedString string) (string, error) {
	lines := strings.Split(encodedString, "\n")
	var decodedLines []string

	for _, line := range lines {
		decodedLine, err := DecodeArt(line)
		if err != nil {
			// In case of any error, including malformed input, return immediately.
			return "", err
		}
		decodedLines = append(decodedLines, decodedLine)
	}

	return strings.Join(decodedLines, "\n"), nil
}
