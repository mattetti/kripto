package kripto

// CharMapScorer is an interface allowing more generic data processing APIs
// while still giving flexibility in how the char maps are being scored.
type CharMapScorer interface {
	Score(m *CharUseMap) float64
}

// EnglishScorer score a character map based on English language statistics.
type EnglishScorer struct {
	// WithSpace indicates if the analyzed text is expected to have white spaces or not.
	WithSpace bool
}

// Score implements the CharMapScorer interface
func (s *EnglishScorer) Score(m *CharUseMap) float64 {
	return m.EnglishScore(s.WithSpace)
}
