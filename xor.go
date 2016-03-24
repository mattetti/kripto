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

// MultiCharXor xors a slice of bytes using a multiple character repeating key
// This is also known as the Vigenère cipher
// https://en.wikipedia.org/wiki/Vigen%C3%A8re_cipher
func MultiCharXor(encStr []byte, k []byte) []byte {
	out := make([]byte, len(encStr))
	kIDX := 0

	for i, b := range encStr {
		out[i] = b ^ k[kIDX]
		kIDX++
		if kIDX > len(k)-1 {
			kIDX = 0
		}
	}
	return out
}

// MostLikelyXorKey tries to guess what the single char xor key might be for a given string.
// Note that this code doesn't do any language analyze and takes a naive/brute force approach.
// Each key is tested against the cypher text and the result map is analyzed.
func MostLikelyXorKey(cypherBlock []byte) byte {
	bestScore := 0.0
	var winnerK byte
	for k := 0; k < 255; k++ {
		data := SingleCharXor(cypherBlock, byte(k))
		cMap := NewCharMap(data)
		score := cMap.ASCIIScore()
		if score >= bestScore {
			bestScore = score
			winnerK = byte(k)
		}
	}
	return winnerK
}

// GuessMultiCharXorKeySize returns a sorted list of most likely key sizes
// for a multi character xor key based on the passed encoded string.
// The slice of possible keys is 10 or less.
func GuessMultiCharXorKeySize(encodedData []byte, maxSize int) []int {
	if len(encodedData) < maxSize {
		maxSize = len(encodedData)
	}
	diffs := diffCol{}
	// for each key size between 2 and max size, we are
	// calculating the hamming distance between blocks of the key size.
	// The key size with the smallest normalized edit distance is probably the key
	for i := 2; i < maxSize; i++ {
		if i*4 > len(encodedData) {
			break
		}
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
