package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Airport represents airports with Name, Age, and Location fields
type Airport struct {
	Name         string `json:"Name"`
	Iso_country  string `json:"Age"`
	Municipality string `json:"Municipality"`
	Icao_code    string `json:"Icao_code"`
	Iata_code    string `json:"Iata_code"`
	Coordinates  string `json:"Coordinates"`
}

// main function is the entry point of the program
func main() {
	helpFlag := false
	flag.BoolVar(&helpFlag, "h", false, "display usage information")
	flag.Parse()

	// Display usage information if the help flag is provided or if the number of command-line arguments is incorrect
	if helpFlag || len(os.Args) != 4 {
		displayUsage()
		return
	}

	// Extract input, output, and airport lookup file paths from command-line arguments
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	airportFile := os.Args[3]

	// Check if input and output files are the same
	if inputFile == outputFile {
		color.Red("output and input file cannot be the same")
		return
	}

	// Open and read the airport lookup CSV file
	file, err := os.Open(airportFile)
	if err != nil {
		color.Red("airport lookup not found")
		return
	}
	defer file.Close()

	///// CSV Reader
	reader := csv.NewReader(file)

	// Read the header of the CSV file
	header, err := reader.Read()
	if err != nil {
		color.Red("error reading CSV header:", err)
		return
	}

	// Check if the CSV file contains all the required columns
	requiredColumns := []string{"name", "icao_code", "municipality", "iso_country", "iata_code", "coordinates"}
	if !hasColumns(header, requiredColumns...) {
		color.Red("CSV does not contain all the required columns.")
		return
	}

	// Read all records from the CSV file and convert each to an Airport struct
	records, _ := reader.ReadAll()
	var airports []Airport
	for _, record := range records {
		airport := Airport{
			Name:         record[0],
			Iso_country:  record[1],
			Municipality: record[2],
			Icao_code:    "##" + record[3],
			Iata_code:    "#" + record[4],
			Coordinates:  record[5],
		}
		airports = append(airports, airport)
	}

	// Read input sentences from the input file
	input, err := checkInput(inputFile)
	// If an error occurs, stop the program
	if err != nil {
		fmt.Println(err)
	} else {
		// Process input, replacing airport tags and formatting dates
		changedSentences, err := processInput(input, airports)
		if err != nil {
			fmt.Println(err)
		} else {
			// Print the changed sentences to the output file
			printResults(changedSentences, outputFile)
		}
	}
}

// displayUsage prints usage information for the command-line tool
func displayUsage() {
	color.Green("itinerary usage:\n")
	color.Red("go run . ./input.txt ./output.txt ./airport-lookup.csv\n")
}

// hasColumns checks if a CSV header contains all the required columns
func hasColumns(header []string, columns ...string) bool {
	for _, colToCheck := range columns {
		found := false
		for _, col := range header {
			if col == colToCheck {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// //
// processInput processes input sentences, replacing airport tags and formatting dates
func processInput(input []string, airports []Airport) ([][]string, error) {
	var changedSentences [][]string

	for _, line := range input {
		// Separate changed sentences from checkTag and testDate
		changedTags, err := checkTag([]string{line}, airports)
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		changedTags = testDate(changedTags)

		changedLine := changedTags
		// Add the changed line to the 2D slice
		changedSentences = append(changedSentences, changedLine)
	}

	return changedSentences, nil
}

// checkInput reads the input file and splits it into a slice of sentences
func checkInput(inputFile string) ([]string, error) {
	inputfile, err := os.Open(inputFile)
	if err != nil {
		return nil, errors.New("input not found")
	}
	defer inputfile.Close()

	var result []string
	scanner := bufio.NewScanner(inputfile)
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line into sentences based on '\n'
		line = strings.ReplaceAll(line, "\\v", "\\n")
		line = strings.ReplaceAll(line, "\\f", "\\n")
		line = strings.ReplaceAll(line, "\\r", "\\n")
		sentences := strings.Split(line, "\\n")
		for _, sentence := range sentences {
			sentence = strings.TrimSpace(sentence)
			result = append(result, sentence)
		}
	}

	if err := scanner.Err(); err != nil {
		errorMessage := "error reading input file"
		color.Red(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return result, nil
}

func checkTag(input []string, airports []Airport) ([]string, error) {
	var changedWords []string

	if isAirportLookupMalformed(airports) {
		errorMessage := "airport lookup malformed"
		color.Red(errorMessage)
		return nil, errors.New(errorMessage)
	}

	for _, words := range input {
		replaceCity := strings.HasPrefix(words, "*")

		for _, airport := range airports {
			// Replace IATA code with airport name if it is a whole word
			if replaceCity && isWholeWord(words, airport.Iata_code) {
				words = strings.Replace(words, airport.Iata_code, airport.Municipality, -1)
				words = strings.TrimPrefix(words, "*")
			}
			if isWholeWord(words, airport.Iata_code) {
				words = strings.Replace(words, airport.Iata_code, airport.Name, -1)
			}
			if replaceCity && isWholeWord(words, airport.Icao_code) {
				words = strings.Replace(words, airport.Icao_code, airport.Municipality, -1)
				words = strings.TrimPrefix(words, "*")
			}
			// Replace ICAO code with airport name if it is a whole word
			if isWholeWord(words, airport.Icao_code) {
				words = strings.Replace(words, airport.Icao_code, airport.Name, -1)
			}
		}

		changedWords = append(changedWords, words)
	}

	return changedWords, nil
}

// isAirportLookupMalformed checks if airport lookup data is malformed
func isAirportLookupMalformed(airports []Airport) bool {
	for _, airport := range airports {
		if airport.Iata_code == "" || airport.Icao_code == "" || airport.Name == "" || airport.Municipality == "" || airport.Iso_country == "" || airport.Coordinates == "" {
			return true
		}
	}
	return false
}

// isWholeWord checks if the given word is a whole word in the text
func isWholeWord(text, word string) bool {
	index := strings.Index(text, word)
	for index != -1 {
		if (index == 0 || !isLetter(text[index-1])) && (index+len(word) == len(text) || !isLetter(text[index+len(word)])) {
			return true
		}
		index = strings.Index(text[index+1:], word)
		if index != -1 {
			index++
		}
	}
	return false
}

// isLetter checks if the character is a letter
func isLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}
