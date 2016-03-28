package kripto

import "fmt"

// ByteKeyStats tracks the stats of a single character encoding key
type ByteKeyStats struct {
	CharMap *CharUseMap
	Score   float64
	Text    []byte
	Key     byte
}

type ByteKeyColStats []*ByteKeyStats

// Len implements the sort interface
func (s ByteKeyColStats) Len() int {
	return len(s)
}

// Swap implements the sort interface
func (s ByteKeyColStats) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implements the sort interface
func (s ByteKeyColStats) Less(i, j int) bool {
	return s[i].Score > s[j].Score
}

func (stats ByteKeyColStats) String() string {
	var o string
	for _, s := range stats {
		o += fmt.Sprintf("%s -> %0.2f -> %s\n", string(s.Key), s.Score, string(s.Text))
	}
	return o
}
