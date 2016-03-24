package kripto

import (
	"bytes"
	"fmt"
	"math"
)

// Debug if true, more print statements will be outputted
var Debug = false

// EnglishLetterFreqs represent the verage % of time each letter appears in English text
var EnglishLetterFreqs = &CharUseMap{
	'a': {Freq: 8.167},
	'b': {Freq: 1.492},
	'c': {Freq: 2.782},
	'd': {Freq: 4.253},
	'e': {Freq: 12.702},
	'f': {Freq: 2.228},
	'g': {Freq: 2.015},
	'h': {Freq: 6.094},
	'i': {Freq: 6.966},
	'j': {Freq: 0.153},
	'k': {Freq: 0.772},
	'l': {Freq: 4.025},
	'm': {Freq: 2.406},
	'n': {Freq: 6.749},
	'o': {Freq: 7.507},
	'p': {Freq: 1.929},
	'q': {Freq: 0.095},
	'r': {Freq: 5.987},
	's': {Freq: 6.327},
	't': {Freq: 9.056},
	'u': {Freq: 2.758},
	'v': {Freq: 0.978},
	'w': {Freq: 2.360},
	'x': {Freq: 0.150},
	'y': {Freq: 1.974},
	'z': {Freq: 0.074},
}

// CharStats represents the metrics of a character from a char map
type CharStats struct {
	Freq  float64
	Count float64
}

// CharUseMap map of usage per character (byte)
type CharUseMap map[byte]*CharStats

// EnglishScore checks if the map looks English.
// WIP
func (m CharUseMap) EnglishScore() (score float64) {
	var penalty float64
	freqs := *EnglishLetterFreqs

	var totalCharCount float64
	for _, charStats := range m {
		totalCharCount += charStats.Count
	}

	maxScore := totalCharCount * 50.0

	for b, charStats := range m {
		engStats, ok := freqs[b]
		if ok {
			penalty += charStats.Count * math.Abs(engStats.Freq-charStats.Freq)
		} else {
			// printable but non letter characters
			if b > 32 && b <= 126 {
				// limited penalty
				penalty += charStats.Count * 5
			} else {
				if IsNumber(b) {
					penalty += charStats.Count * 2
					continue
				}
				if IsSpace(b) {
					// average word length 5.1 characters
					expectedSpaces := (totalCharCount / 5.1) - 1
					expF := (expectedSpaces * 100) / totalCharCount
					minExpF := (((totalCharCount / 3) - 1) * 100) / totalCharCount

					if charStats.Freq < minExpF {
						penalty += totalCharCount * 100
						continue
					}

					penalty += charStats.Count * math.Abs(expF-charStats.Freq)
					continue
				}
				// highest penalty for non printable chars
				penalty += charStats.Count * 50
			}
		}
	}

	return maxScore - penalty
}

// ASCIIScore allocates a score based on ASCII like the char map is.
func (m *CharUseMap) ASCIIScore() float64 {
	var penalty float64

	var totalCharCount float64
	for _, charStats := range *m {
		totalCharCount += charStats.Count
	}
	maxScore := totalCharCount * 50.0

	for b, charStats := range *m {
		if !IsPrintable(b) {
			penalty += charStats.Count * 100
			continue
		}
		// treating letters and numbers are what we want
		if IsASCIILetter(b) || IsNumber(b) {
			continue
		}
		if IsSpace(b) {
			if charStats.Freq > 0.3 {
				penalty += charStats.Count * 30
			}
			continue
		}
		if IsPunctuation(b) {
			// check the frequency
			if charStats.Freq < 0.1 {
				// probably fine
				penalty += charStats.Count * 10
				continue
			}
			if charStats.Freq < 0.3 {
				penalty += charStats.Count * 10
				continue
			}
			penalty += charStats.Count * 40
			continue
		}
		penalty += charStats.Count * 60
	}

	return maxScore - penalty
}

func (m *CharUseMap) String() string {
	var out string
	for b, stats := range *m {
		out += fmt.Sprintf("%s: %+v\n", string(b), *stats)
	}
	return out
}

// IsPunctuation checks if the provided byte can be considered as punctuation
func IsPunctuation(b byte) bool {
	switch b {
	// !"'(),:;?
	case 33, 34, 39, 40, 41, 44, 46, 58, 59, 63:
		return true
	default:
		return false
	}
}

// IsLetter checks if the byte is an ASCII letter (upper or lower case)
func IsASCIILetter(b byte) bool {
	// uppercase letter
	if b >= 65 && b <= 90 {
		return true
	}
	// lowercase letter
	if b >= 97 && b <= 122 {
		return true
	}
	return false
}

// IsNumber indicates that the byte represents a number or not
func IsNumber(b byte) bool {
	if b >= 48 && b <= 57 {
		return true
	}
	return false
}

// IsSpace does what you expect it to do
func IsSpace(b byte) bool {
	return b == 32
}

func IsPrintable(b byte) bool {
	return b >= 32 && b <= 126
}

// CharMaps are a colletion of CharUseMaps so we can compare them to each other.
type CharMaps []*CharUseMap

// MostEnglish returns the most English charmap based on score
func (maps *CharMaps) MostEnglish() *CharUseMap {
	bestScore := 0.0
	winnerIDX := -1
	for i, cmap := range *maps {
		score := cmap.EnglishScore()
		if score > bestScore {
			winnerIDX = i
			bestScore = score
		}
		if Debug {
			fmt.Printf("%d - %f\n", i, score)
		}
	}
	if winnerIDX == -1 {
		return nil
	}
	m := *maps
	return m[winnerIDX]
}

// NewCharMap returns a map of the character usage (converted to lowercase)
func NewCharMap(str []byte) *CharUseMap {
	charCount := map[byte]int{}
	for _, b := range str {
		// lowercase ASCII letters
		if IsASCIILetter(b) {
			b = bytes.ToLower([]byte{b})[0]
		}
		charCount[b]++
	}
	// TODO: word lengths
	m := CharUseMap{' ': &CharStats{}}
	for b, count := range charCount {
		m[b] = &CharStats{
			Count: float64(count),
			Freq:  (100.0 * float64(count)) / float64(len(charCount)),
		}
	}
	return &m
}
