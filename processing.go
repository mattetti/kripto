package kripto

import (
	"encoding/base64"
	"encoding/hex"
)

// DeHex decodes a hex encoded string (but swallows errors)
func DeHex(input []byte) []byte {
	decStr := make([]byte, hex.DecodedLen(len(input)))
	hex.Decode(decStr, input)
	return decStr
}

// DeBase64 decodes a base64 string
func DeBase64(input []byte) []byte {
	enc := base64.StdEncoding
	dbuf := make([]byte, enc.DecodedLen(len(input)))
	n, err := enc.Decode(dbuf, input)
	if err != nil {
		panic(err)
	}
	return dbuf[:n]
}
