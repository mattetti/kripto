package kripto

// diffCol is a collection of diffs
type diffCol []*diff

// diff tracks the normalized difference for a given value
type diff struct {
	// the analyzed val
	val int
	// the normalized difference
	normDiff float64
}

// Len implements the sort interface
func (d diffCol) Len() int {
	return len(d)
}

// Swap implements the sort interface
func (d diffCol) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less implements the sort interface
func (d diffCol) Less(i, j int) bool {
	return d[i].normDiff < d[j].normDiff
}
