package kripto

// SingleCharXor xors a slice of bytes using the passed key
func SingleCharXor(encStr []byte, k byte) []byte {
	outB := make([]byte, len(encStr))
	for i, b := range encStr {
		outB[i] = b ^ k
	}
	return outB
}
