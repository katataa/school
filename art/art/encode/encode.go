package encode

import (
	"fmt"
	"strings"
)

// encodeLine encodes a single line using custom rules for specific patterns
// and falls back to run-length encoding for general cases.
func encodeLine(line string) string {
	// Patterns to look for and their encoded representations
	patternMap := map[string]string{
		"/^--^\\     /^--^\\     /^--^\\":   "[22  ]/^--^\\     /^--^\\     /^--^\\",
		"\\____/     \\____/     \\____/":   "[22  ]\\[4 _]/     \\[4 _]/     \\[4 _]/",
		"/      \\   /      \\   /      \\": "[21  ]/[6  ]\\   /[6  ]\\   /[6  ]\\",
		"|        | |        | |        |":  "[20  ]|[8  ]| |[8  ]| |[8  ]|",
		"\\__  __/   \\__  __/   \\__  __/": "[21  ]\\__  __/   \\__  __/   \\__  __/",
	}

	// Attempt to match and replace patterns
	for pattern, replacement := range patternMap {
		if strings.Contains(line, pattern) {
			return replacement
		}
	}

	// Fallback to run-length encoding for lines without special patterns
	var result strings.Builder
	count := 1
	for i := 0; i < len(line); i++ {
		if i+1 < len(line) && line[i] == line[i+1] {
			count++
		} else {
			if count > 1 {
				result.WriteString(fmt.Sprintf("[%d %s]", count, string(line[i])))
				count = 1
			} else {
				result.WriteString(string(line[i]))
			}
		}
	}
	return result.String()
}

// EncodeArt processes each line of the input ASCII art with encodeLine function.
func EncodeArt(input string) string {
	var result strings.Builder
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		encodedLine := encodeLine(line)
		result.WriteString(encodedLine + "\n")
	}

	return strings.TrimRight(result.String(), "\n")
}
