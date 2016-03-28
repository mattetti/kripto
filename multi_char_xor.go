package kripto

// BreakMultiCharXor tries to break encoded text that was encrypted using a repeating multiple character key
// XORing the data.
func BreakMultiCharXor(data []byte, maxKeyLength int) (out []byte, k []byte) {
	// possible key sizes
	kSizes := GuessMultiCharXorKeySize(data, maxKeyLength)
	// test the first possible key (TODO: might want to test more)
	max := 1
	if len(kSizes) < max {
		max = len(kSizes)
	}
	possibleKeys := make([][]byte, max)

	for idx, kSize := range kSizes[:max] {
		// break the ciphertext into blocks of key size length
		blocks := [][]byte{}
		for i := 0; i+kSize <= len(data); i = i + kSize {
			blocks = append(blocks, data[i:i+kSize])
		}

		// fmt.Printf("%#v\n", blocks)

		// transpose the blocks: make a block that is the first byte of every block,
		// and a block that is the second byte of every block, and so on.
		// Each extracted block contains characters encoded with a single character xor key.
		xordBlocks := make([][]byte, kSize)
		for _, block := range blocks {
			for i, b := range block {
				if len(xordBlocks[i]) == 0 {
					xordBlocks[i] = []byte{}
				}
				xordBlocks[i] = append(xordBlocks[i], b)
			}
		}

		key := make([]byte, kSize)
		for i, cypherBlock := range xordBlocks {
			// fmt.Printf("%d %#v\n", i, cypherBlock)
			key[i] = MostLikelyXorKey(cypherBlock)
		}
		possibleKeys[idx] = key
	}

	return MultiCharXor(data, possibleKeys[0]), possibleKeys[0]
}
