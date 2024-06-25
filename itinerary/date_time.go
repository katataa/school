package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

///date and time conversion :)

// testDate formats dates in different ways based on prefixes in the input
func testDate(input []string) []string {
	var changedDates []string

	for _, inputDate := range input {
		var layout string
		var formattedDate string

		switch {
		case strings.HasPrefix(inputDate, "D("):
			layout = "D(2006-01-02T15:04Z07:00)"
			parsedDate, err := time.Parse(layout, inputDate)
			if err != nil {
				fmt.Println("error parsing date:", err)
			}
			formattedDate = parsedDate.Format("02 Jan 2006")
		case strings.HasPrefix(inputDate, "T12("):
			layout = "T15(2006-01-02T15:04Z07:00)"
			parsedDate, err := time.Parse(layout, inputDate)
			if err != nil {
				// If parsing fails, leave the date unchanged
				formattedDate = inputDate
				break
			}

			if strings.HasSuffix(inputDate, "Z") {
				formattedDate = parsedDate.Format("03:04PM (+00:00)")
			} else {
				formattedDate = parsedDate.Format("03:04PM (-07:00)")
			}
		case strings.HasPrefix(inputDate, "T24("):
			// Extract the date-time substring
			dateString := inputDate[4 : len(inputDate)-1]

			var parsedDate time.Time
			var err error
			if strings.HasSuffix(dateString, "Z") {
				parsedDate, err = time.Parse("2006-01-02T15:04", dateString[:len(dateString)-1])
			} else {
				parsedDate, err = time.Parse("2006-01-02T15:04-07:00", dateString)
			}

			if err != nil {
				fmt.Println("error parsing date:", err)
			}

			// Format the parsed date
			if strings.HasSuffix(dateString, "Z") {
				formattedDate = parsedDate.Format("15:04 (+00:00)")
			} else {
				formattedDate = parsedDate.Format("15:04 (-07:00)")
			}

		default:
			formattedDate = inputDate
		}

		changedDates = append(changedDates, formattedDate)
	}

	return changedDates
}

// printResults writes the changed sentences to an output file
func printResults(changedSentences [][]string, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	prevLineEmpty := false

	// Print each line separately with a newline character
	for _, line := range changedSentences {
		currentLine := strings.Join(line, " ") + "\n"

		isEmpty := len(strings.TrimSpace(currentLine)) == 0

		if isEmpty && prevLineEmpty {
			continue
		}

		// Use color formatting for specific information (modify as needed)
		currentLine = highlightInformation(currentLine)

		_, err = file.WriteString(currentLine)
		if err != nil {
			return err
		}

		prevLineEmpty = isEmpty
	}

	return nil
}

// highlightInformation applies color formatting to specific information
func highlightInformation(line string) string {
	// Example: Highlight dates in green
	line = strings.ReplaceAll(line, "05 Apr 2007", color.GreenString("05 Apr 2007"))

	// Example: Highlight times in blue
	line = strings.ReplaceAll(line, "12:30PM (-02:00)", color.BlueString("12:30PM (-02:00)"))

	// Add more formatting as needed

	return line
}
