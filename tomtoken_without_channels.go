package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	// Parse command line arguments
	outputFlag := flag.String("output", "stdout", "Output: 'stdout' or path to .json file")
	tokensFlag := flag.Int("tokens", 100_000_000, "Number of TomToken tokens to generate")
	flag.Parse()

	// Setup output (stdout or file)
	var outFile *os.File
	if *outputFlag == "stdout" {
		outFile = os.Stdout
	} else {
		f, err := os.Create(*outputFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		outFile = f
	}
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	// Write metadata
	created := time.Now().UTC().Format(time.RFC3339)
	_, _ = writer.WriteString("{\"metadata\":{\"title\":\"Generated TomToken Stream\",\"created\":\"")
	_, _ = writer.WriteString(created)
	_, _ = writer.WriteString("\"}, \"referenceMap\":{")

	// Write referenceMap table (space identifiers)
	// Define ID and character pairs:
	refEntries := []struct {
		id    string
		value string
	}{
		{"1", " "},     // space
		{"2", "\t"},    // tab
		{"3", "\u00A0"}, // non-breaking space (NBSP)
	}
	for i, entry := range refEntries {
		var escaped string
		// Need to escape control characters for JSON
		if entry.value == "\t" {
			escaped = "\\t"
		} else if entry.value == "\u00A0" {
			escaped = "\\u00a0"
		} else {
			escaped = entry.value
		}
		_, _ = writer.WriteString("\"" + entry.id + "\": \"" + escaped + "\"")
		if i < len(refEntries)-1 {
			_, _ = writer.WriteString(", ")
		}
	}
	_, _ = writer.WriteString("}, \"content\":[\n")

	// Initialize random number generator
	rand.Seed(time.Now().UnixNano())

	// Variables for line length control and state
	lineLen := 0      // current line length in characters
	firstToken := true
	total := *tokensFlag

	// Generate and stream content
	for generated := 0; generated < total; {
		// 1. Generate random word 3-12 letters
		wordLen := 3 + rand.Intn(10) // random length from 3 to 12
		wordBytes := make([]byte, wordLen)
		for i := 0; i < wordLen; i++ {
			wordBytes[i] = byte('a' + rand.Intn(26))
		}
		token := "\"" + string(wordBytes) + "\"" // string token in quotes

		// 2. Write word to JSON (considering first element and line breaks if needed)
		if firstToken {
			// First array element (line break already done after "[")
			if lineLen + len(token) > 250 {
				_, _ = writer.WriteString("\n")
				lineLen = 0
			}
			_, _ = writer.WriteString(token)
			lineLen += len(token)
			firstToken = false
		} else {
			// Not the first element - need a comma before token
			if lineLen + 1 + len(token) > 250 {
				// Doesn't fit on current line - move to new line
				_, _ = writer.WriteString(",\n")
				lineLen = 0
				_, _ = writer.WriteString(token)
				lineLen += len(token)
			} else {
				// Fits on the same line
				_, _ = writer.WriteString(", " + token)
				lineLen += 2 + len(token)
			}
		}
		generated++
		if generated >= total {
			break
		}

		// 3. Decide what to insert after the word: space or line break?
		if rand.Float64() < 0.01 {
			// **Insert line break** (rare)
			modes := []string{"lf", "cr", "crlf"}
			mode := modes[rand.Intn(len(modes))]
			newlineToken := "{\"type\":\"newline\",\"mode\":\"" + mode + "\"}"
			// With 10% probability make an empty line (duplicate line break)
			if rand.Float64() < 0.1 {
				newlineToken = "{\"type\":\"repeat\",\"count\":2,\"value\":" + newlineToken + "}"
			}
			// Write line break token
			if lineLen + 1 + len(newlineToken) > 250 {
				_, _ = writer.WriteString(",\n")
				lineLen = 0
				_, _ = writer.WriteString(newlineToken)
				lineLen += len(newlineToken)
			} else {
				_, _ = writer.WriteString(", " + newlineToken)
				lineLen += 2 + len(newlineToken)
			}
			generated++
			if generated >= total {
				break
			}

			// 3a. After line break, sometimes insert indentation (tab or spaces)
			if rand.Float64() < 0.3 {
				var indentToken string
				if rand.Float64() < 0.5 {
					indentToken = "\"2\"" // tab (reference)
				} else {
					indentToken = "{\"type\":\"repeat\",\"count\":4,\"value\":1}" // 4 spaces as one token
				}
				if lineLen + 1 + len(indentToken) > 250 {
					_, _ = writer.WriteString(",\n")
					lineLen = 0
					_, _ = writer.WriteString(indentToken)
					lineLen += len(indentToken)
				} else {
					_, _ = writer.WriteString(", " + indentToken)
					lineLen += 2 + len(indentToken)
				}
				generated++
			}
		} else {
			// **Insert space** (common case)
			var spaceToken string
			x := rand.Float64()
			if x < 0.02 {
				// Double or triple space
				repeatCount := 2 + rand.Intn(3) // 2-4
				spaceToken = fmt.Sprintf("{\"type\":\"repeat\",\"count\":%d,\"value\":1}", repeatCount)
			} else if x < 0.07 {
				// Non-breaking space instead of regular
				spaceToken = "\"3\""
			} else {
				// Regular space
				spaceToken = "\"1\""
			}

			if lineLen + 1 + len(spaceToken) > 250 {
				_, _ = writer.WriteString(",\n")
				lineLen = 0
				_, _ = writer.WriteString(spaceToken)
				lineLen += len(spaceToken)
			} else {
				_, _ = writer.WriteString(", " + spaceToken)
				lineLen += 2 + len(spaceToken)
			}
			generated++
		}
		if generated >= total {
			break
		}
		// After space, cycle continues with next word generation
	}

	// 4. Close JSON array and object, considering last line length
	if lineLen + 2 > 250 {
		_, _ = writer.WriteString("\n")
	}
	_, _ = writer.WriteString("]}")
}
