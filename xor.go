package kripto

import "sort"

// DataProcessFn offers a generic interface to process input data
type DataProcessFn func(data []byte) []byte

// SingleCharXor xors a slice of bytes using the passed key
func SingleCharXor(encStr []byte, k byte) []byte {
	outB := make([]byte, len(encStr))
	for i, b := range encStr {
		outB[i] = b ^ k
	}
	return outB
}

// BreakSingleCharXor does its best to break English text encrypted using a single character xor key.
// The processFn param is a function that can be used to apply basic input processing (hex/base64 decoding for instance).
// The scorer param is used to score the output data and find the right key.
func BreakSingleCharXor(xord []byte, processFn DataProcessFn, scorer CharMapScorer) (out []byte, key byte) {
	if processFn != nil {
		xord = processFn(xord)
	}

	stats := ByteKeyColStats{}
	for k := byte(0); k < 255; k++ {
		text := SingleCharXor(xord, k)
		m := NewCharMap(text)
		stats = append(stats, &ByteKeyStats{
			CharMap: m,
			Score:   scorer.Score(m),
			Text:    text,
			Key:     k,
		})
	}

	if len(stats) == 0 {
		return
	}

	sort.Sort(stats)
	return stats[0].Text, stats[0].Key
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
// This is a naive brute force approach that looks at the output text and checks its statistical validity
// as English text (using spaces).
func MostLikelyXorKey(cypherBlock []byte) byte {
	bestScore := 0.0
	var winnerK byte
	for k := 0; k < 255; k++ {
		data := SingleCharXor(cypherBlock, byte(k))
		cMap := NewCharMap(data)
		score := cMap.EnglishScore(true)
		if score >= bestScore {
			// give a preference to ascii letters
			if score == bestScore && IsASCIILetter(winnerK) {
				continue
			}
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

	for keySize := 2; keySize < maxSize; keySize++ {
		chunks := [][]byte{}
		// extract up to 4 chunks of data
		for i := 0; i < 4; i++ {
			if keySize*(i+1) > len(encodedData) {
				break
			}
			chunks = append(chunks, encodedData[keySize*i:keySize*(i+1)])
		}

		diffSum := 0.0
		for i, c := range chunks {
			for _, otherC := range chunks[i:] {
				diffSum += float64(HammingDiff(c, otherC))
			}
		}

		diffs = append(diffs, &diff{
			val:      keySize,
			normDiff: diffSum / float64(len(chunks)) / float64(keySize),
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
