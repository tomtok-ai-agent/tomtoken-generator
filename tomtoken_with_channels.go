package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"math/rand"
	"os"
	"strings"
	"time"
)

// TomToken represents a token type
type TomToken struct {
	Type  string      `json:"type,omitempty"`
	Mode  string      `json:"mode,omitempty"`
	Count int         `json:"count,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// Settings for the generator channel
type GeneratorSettings struct {
	TokenCount int
}

// Token generator using channels
func GenerateTokens(settings GeneratorSettings) <-chan interface{} {
	ch := make(chan interface{}, 100) // buffered channel for efficiency

	go func() {
		defer close(ch)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		// Space token IDs: 1 = regular space, 3 = non-breaking space

		for generated := 0; generated < settings.TokenCount; {
			// Generate a word
			wordLen := 3 + r.Intn(10)
			wordBuilder := strings.Builder{}
			for i := 0; i < wordLen; i++ {
				wordBuilder.WriteByte(byte('a' + r.Intn(26)))
			}
			ch <- wordBuilder.String()
			generated++
			if generated >= settings.TokenCount {
				break
			}

			// What to insert after the word?
			p := r.Float64()
			if p < 0.01 {
				// Newline
				newlineModes := []string{"lf", "crlf", "cr"}
				mode := newlineModes[r.Intn(len(newlineModes))]
				token := TomToken{Type: "newline", Mode: mode}

				// Sometimes make a repeat
				if r.Float64() < 0.1 {
					token = TomToken{
						Type:  "repeat",
						Count: 2,
						Value: token,
					}
				}
				ch <- token
				generated++
				if generated >= settings.TokenCount {
					break
				}

				// Sometimes insert indentation
				if r.Float64() < 0.3 {
					if r.Float64() < 0.5 {
						ch <- 2 // tab
					} else {
						ch <- TomToken{Type: "repeat", Count: 4, Value: 1} // 4 spaces
					}
					generated++
				}
			} else {
				// Spaces and repeats
				if p < 0.02 {
					ch <- TomToken{Type: "repeat", Count: 2 + r.Intn(3), Value: 1}
				} else if p < 0.07 {
					ch <- 3 // non-breaking space
				} else {
					ch <- 1 // regular space
				}
				generated++
			}
		}
	}()

	return ch
}

func main() {
	// Command line parameters
	outputFile := flag.String("output", "stdout", "output: stdout or file")
	tokenCount := flag.Int("tokens", 100_000, "number of tokens")
	flag.Parse()

	var out *os.File
	var err error

	if *outputFile == "stdout" {
		out = os.Stdout
	} else {
		out, err = os.Create(*outputFile)
		if err != nil {
			panic(err)
		}
		defer out.Close()
	}

	writer := bufio.NewWriter(out)
	defer writer.Flush()

	// Start forming JSON
	header := map[string]interface{}{
		"metadata": map[string]string{
			"title":   "Generated TomToken Stream",
			"created": time.Now().UTC().Format(time.RFC3339),
		},
		"referenceMap": map[string]string{
			"1": " ",
			"2": "\t",
			"3": "\u00a0",
		},
	}

	encoder := json.NewEncoder(writer)
	writer.WriteString("{\n")

	// Write metadata and referenceMap
	writer.WriteString(`"metadata": `)
	encoder.Encode(header["metadata"])
	writer.WriteString(",\n")

	writer.WriteString(`"referenceMap": `)
	encoder.Encode(header["referenceMap"])
	writer.WriteString(",\n")

	writer.WriteString(`"content": [` + "\n")

	// Stream generation and writing of tokens
	tokenCh := GenerateTokens(GeneratorSettings{TokenCount: *tokenCount})
	first := true
	lineLen := 0

	for token := range tokenCh {
		var tokenData []byte
		switch v := token.(type) {
		case string, int:
			tokenData, _ = json.Marshal(v)
		case TomToken:
			tokenData, _ = json.Marshal(v)
		default:
			continue
		}

		// Control line length (250 characters)
		tokenStr := string(tokenData)
		if !first {
			tokenStr = "," + tokenStr
		}
		if lineLen+len(tokenStr) > 250 {
			writer.WriteString("\n")
			lineLen = 0
			tokenStr = strings.TrimPrefix(tokenStr, ",") // comma is already on the previous line
		}

		writer.WriteString(tokenStr)
		lineLen += len(tokenStr)
		first = false
	}

	// Complete JSON
	writer.WriteString("\n]}\n")
}
