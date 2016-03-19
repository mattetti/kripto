package kripto

import "sort"

// SingleCharXor xors a slice of bytes using the passed key
func SingleCharXor(encStr []byte, k byte) []byte {
	outB := make([]byte, len(encStr))
	for i, b := range encStr {
		outB[i] = b ^ k
	}
	return outB
}

// GuessMultiCharXorKeySize returns a sorted list of most likely key sizes
// for a multi character xor key based on the passed encoded string.
// The slice of possible keys is 10 or less.
func GuessMultiCharXorKeySize(encodedData []byte, maxSize int) []int {
	diffs := diffCol{}
	// for each key size between 2 and max size, we are
	// calculating the hamming distance between blocks of the key size.
	// The key size with the smallest normalized edit distance is probably the key
	for i := 2; i < maxSize+1; i++ {
		first := float64(HammingDiff(encodedData[:i], encodedData[i:i+i])) / float64(i)
		second := float64(HammingDiff(encodedData[i*2:i*3], encodedData[i*3:i*4])) / float64(i)

		diffs = append(diffs, &diff{
			val:      i,
			normDiff: (first + second) / 2,
		})
	}
	sort.Sort(diffs)

	// returning 10 or less result so the consumer can try various keys
	max := 10
	if len(diffs) < max {
		max = len(diffs)
	}
	possibleKeys := make([]int, max)
	for i := 0; i < max; i++ {
		possibleKeys[i] = diffs[i].val
	}
	return possibleKeys
}
