package badwords

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// badWords is a list of bad words or patterns loaded from a text file.
var badWords []string

// LoadBadWords loads bad words from a text file.
// Each line in the file represents a bad word or pattern.
func LoadBadWords(filename string) (bool, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return false, err
	}

	// Split the file content into lines
	badWords = strings.Split(string(data), "\n")

	// Trim whitespace and remove empty lines
	for i := 0; i < len(badWords); i++ {
		badWords[i] = strings.TrimSpace(badWords[i])
		if badWords[i] == "" {
			badWords = append(badWords[:i], badWords[i+1:]...)
			i-- // Adjust index after removing an empty line
		}
	}

	fmt.Printf("Loaded %d bad words from text file\n", len(badWords))
	return true, nil
}

// ContainsBadWords checks if the input text contains any bad words.
func ContainsBadWords(text string) bool {
	lowerText := strings.ToLower(text)
	for _, badWord := range badWords {
		if strings.Contains(lowerText, badWord) {
			fmt.Printf("Bad word detected: %s\n", badWord) // Log the detected bad word
			return true
		}
	}
	return false
}

// AddBadWord adds a new bad word to the list.
func AddBadWord(badWord string) (bool, error) {
	if badWord == "" {
		return false, errors.New("bad word must not be empty")
	}
	badWords = append(badWords, badWord)
	return true, nil
}

// RemoveBadWord removes a bad word from the list.
func RemoveBadWord(badWord string) bool {
	for i, bw := range badWords {
		if bw == badWord {
			badWords = append(badWords[:i], badWords[i+1:]...)
			return true
		}
	}
	return false
}

// ListBadWords returns the current list of bad words.
func ListBadWords() []string {
	return badWords
}