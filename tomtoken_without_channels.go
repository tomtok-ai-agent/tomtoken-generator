package main

import (
	"fmt"
	"math/rand"
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

// Generate tokens without using channels
func tokenGenerator(count int) []interface{} {
	tokens := make([]interface{}, 0, count)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Space token IDs: 1 = regular space, 3 = non-breaking space

	for generated := 0; generated < count; {
		// Generate a word
		wordLen := 3 + r.Intn(10)
		wordBuilder := strings.Builder{}
		for i := 0; i < wordLen; i++ {
			wordBuilder.WriteByte(byte('a' + r.Intn(26)))
		}
		tokens = append(tokens, wordBuilder.String())
		generated++
		if generated >= count {
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
			tokens = append(tokens, token)
			generated++
			if generated >= count {
				break
			}

			// Sometimes insert indentation
			if r.Float64() < 0.3 {
				if r.Float64() < 0.5 {
					tokens = append(tokens, 2) // tab
				} else {
					tokens = append(tokens, TomToken{Type: "repeat", Count: 4, Value: 1}) // 4 spaces
				}
				generated++
			}
		} else {
			// Spaces and repeats
			if p < 0.02 {
				tokens = append(tokens, TomToken{Type: "repeat", Count: 2 + r.Intn(3), Value: 1})
			} else if p < 0.07 {
				tokens = append(tokens, 3) // non-breaking space
			} else {
				tokens = append(tokens, 1) // regular space
			}
			generated++
		}
	}

	return tokens
}

func main() {
	tokens := tokenGenerator(100_000)
	for _, token := range tokens {
		// Process and output token to JSON
		fmt.Println(token)
	}
}
