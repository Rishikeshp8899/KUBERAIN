package utility

import (
	"math"
	"strings"
)

// CosineSimilarity calculates the text similarity between an RFQ text and a supplier's capability summary using term-frequency vectors.
func CosineSimilarity(text1, text2 string) float64 {
	// Tokenize and clean texts
	words1 := tokenize(text1)
	words2 := tokenize(text2)

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	// Build frequency maps
	freq1 := make(map[string]float64)
	for _, w := range words1 {
		freq1[w]++
	}

	freq2 := make(map[string]float64)
	for _, w := range words2 {
		freq2[w]++
	}

	// Calculate dot product and magnitudes
	var dotProduct, mag1, mag2 float64

	allWords := make(map[string]bool)
	for w := range freq1 {
		allWords[w] = true
	}
	for w := range freq2 {
		allWords[w] = true
	}

	for w := range allWords {
		v1 := freq1[w]
		v2 := freq2[w]
		dotProduct += v1 * v2
		mag1 += v1 * v1
		mag2 += v2 * v2
	}

	if mag1 == 0 || mag2 == 0 {
		return 0.0
	}

	similarity := dotProduct / (math.Sqrt(mag1) * math.Sqrt(mag2))
	return math.Round(similarity*1000) / 1000 // rounded to 3 decimal places
}

func tokenize(text string) []string {
	// Simple lowercasing and punctuation stripping split
	cleaned := strings.ToLower(text)
	// Replace punctuation with spaces
	mapping := func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return ' '
	}
	fields := strings.Fields(strings.Map(mapping, cleaned))

	// Optional: filter stop words if necessary
	return fields
}

// ComputeCosineSimilarity calculates text vector similarity (0.0 to 1.0) between two strings.
func ComputeCosineSimilarity(textA, textB string) float64 {
	tokensA := tokenizeText(textA)
	tokensB := tokenizeText(textB)

	if len(tokensA) == 0 || len(tokensB) == 0 {
		return 0.0
	}

	freqA := make(map[string]float64)
	for _, t := range tokensA {
		freqA[t]++
	}

	freqB := make(map[string]float64)
	for _, t := range tokensB {
		freqB[t]++
	}

	var dotProduct, normA, normB float64
	vocab := make(map[string]bool)
	for k := range freqA {
		vocab[k] = true
	}
	for k := range freqB {
		vocab[k] = true
	}

	for word := range vocab {
		vA := freqA[word]
		vB := freqB[word]
		dotProduct += vA * vB
		normA += vA * vA
		normB += vB * vB
	}

	if normA == 0 || normB == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func tokenizeText(input string) []string {
	cleaned := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return ' '
	}, strings.ToLower(input))

	stopWords := map[string]bool{
		"and": true, "the": true, "for": true, "with": true, "of": true, "to": true, "in": true, "a": true, "an": true,
	}

	var tokens []string
	for _, word := range strings.Fields(cleaned) {
		if len(word) > 1 && !stopWords[word] {
			tokens = append(tokens, word)
		}
	}
	return tokens
}
