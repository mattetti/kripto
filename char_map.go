package kripto

import (
	"bytes"
	"fmt"
	"math"
)

// Debug if true, more print statements will be outputted
var Debug = false

// EnglishLetterFreqs represent the verage % of time each letter appears in English text
// From https://en.wikipedia.org/wiki/Letter_frequency#Relative_frequencies_of_letters_in_the_English_language
var EnglishLetterFreqs = &CharUseMap{
	'a': {Freq: 0.08167},
	'b': {Freq: 0.01492},
	'c': {Freq: 0.02782},
	'd': {Freq: 0.04253},
	'e': {Freq: 0.12702},
	'f': {Freq: 0.02228},
	'g': {Freq: 0.02015},
	'h': {Freq: 0.06094},
	'i': {Freq: 0.06966},
	'j': {Freq: 0.00153},
	'k': {Freq: 0.00772},
	'l': {Freq: 0.04025},
	'm': {Freq: 0.02406},
	'n': {Freq: 0.06749},
	'o': {Freq: 0.07507},
	'p': {Freq: 0.01929},
	'q': {Freq: 0.00095},
	'r': {Freq: 0.05987},
	's': {Freq: 0.06327},
	't': {Freq: 0.09056},
	'u': {Freq: 0.02758},
	'v': {Freq: 0.00978},
	'w': {Freq: 0.02360},
	'x': {Freq: 0.00150},
	'y': {Freq: 0.01974},
	'z': {Freq: 0.00074},
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
func (m CharUseMap) EnglishScore(spacedText bool) (score float64) {
	var penalty float64
	var letterPenalty float64
	freqs := *EnglishLetterFreqs

	// expected letters
	for b, en := range freqs {
		stats, ok := m[b]
		if ok {
			letterPenalty += en.Freq * math.Abs(en.Freq-stats.Freq)
		} else {
			letterPenalty += en.Freq * en.Freq
		}
	}

	var spacePenalty float64
	if spacedText {
		stats, ok := m[' ']
		if ok {
			spacePenalty = 0.15 * math.Abs(stats.Freq-0.15)
		} else {
			spacePenalty = 0.15
		}
	}

	var punctCount float64
	var punctFs float64
	for b, stats := range m {
		if IsPunctuation(b) {
			punctCount += stats.Count
			punctFs += stats.Freq
			continue
		}
		if !IsPrintable(b) {
			penalty += math.Abs(stats.Freq + 0.2)
		}
	}
	punctuationPenalty := math.Abs(punctFs - 0.02)

	penalty += (letterPenalty + spacePenalty + punctuationPenalty)
	return math.Max((1.0 - penalty), 0.0)
}

// ASCIIScore allocates a score based on ASCII like the char map is.
func (m *CharUseMap) ASCIIScore() float64 {
	var penalty float64

	var totalCharCount float64
	for _, charStats := range *m {
		totalCharCount += charStats.Count
	}
	score := 0.0

	for b, charStats := range *m {
		if !IsPrintable(b) {
			penalty += charStats.Count * 100
			continue
		}
		// treating letters and numbers are what we want
		if IsASCIILetter(b) || IsNumber(b) {
			score += charStats.Count * 70
		}
		if IsSpace(b) {
			if charStats.Freq > 0.5 {
				penalty += charStats.Count * 30
				continue
			}
			score += charStats.Count * 30
			continue
		}
		if IsPunctuation(b) {
			// check the frequency
			if charStats.Freq < 0.1 {
				// probably fine
				score += charStats.Count * 10
				continue
			}
			if charStats.Freq < 0.3 {
				score += charStats.Count * 5
				continue
			}
			penalty += charStats.Count * 30
			continue
		}
		penalty += charStats.Count * 60
	}

	return score - penalty
}

func (m *CharUseMap) String() string {
	var out string
	for b, stats := range *m {
		out += fmt.Sprintf("%s: %+v\n", string(b), *stats)
	}
	return out
}

// IsPunctuation checks if the provided byte can be considered as "punctuation"
func IsPunctuation(b byte) bool {
	// TODO: maybe should rename and use IsASCIILetter with better bounds to determine, printable
	// non letter/numbers chars.
	switch b {
	case '"', '\'', '(', ')', ',', ':', ';', '.', '?', '!':
		return true
	case '[', ']', '{', '}', '+', '-', '~', '`', '*', '%', '@', '$', '<', '>', '|', '_', '\\', '/', '#', '^':
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

// IsPrintable checks that a byte can be printed
func IsPrintable(b byte) bool {
	switch b {
	case 9, 10, 12, 13:
		return true
	default:
		return (b >= 32 && b <= 126)
	}

}

// CharMaps are a colletion of CharUseMaps so we can compare them to each other.
type CharMaps []*CharUseMap

// MostEnglish returns the most English charmap based on score
func (maps *CharMaps) MostEnglish() *CharUseMap {
	bestScore := 0.0
	winnerIDX := -1
	for i, cmap := range *maps {
		score := cmap.EnglishScore(true)
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
	m := CharUseMap{}
	for b, count := range charCount {
		m[b] = &CharStats{
			Count: float64(count),
			Freq:  float64(count) / float64(len(charCount)),
		}
	}
	return &m
}
