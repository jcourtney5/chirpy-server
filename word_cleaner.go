package main

import (
	"strings"
)

func wordCleaner(sentence string) string {
	// Create a set of profane words
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	const replacement = "****"

	words := strings.Split(sentence, " ")
	cleanedWords := make([]string, 0, len(words))
	for _, word := range words {
		wordLower := strings.ToLower(word)
		if _, exists := profaneWords[wordLower]; exists {
			cleanedWords = append(cleanedWords, replacement)
		} else {
			cleanedWords = append(cleanedWords, word)
		}
	}

	return strings.Join(cleanedWords, " ")
}
