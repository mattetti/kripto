package kripto

import (
	"crypto/aes"
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"
)

func TestECBDecrypter(t *testing.T) {
	key := []byte("YELLOW SUBMARINE")
	f, err := os.Open("fixtures/7.txt")
	if err != nil {
		t.Fatal(err)
	}
	ciphertext, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	ciphertext, err = base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		t.Fatal(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}

	// ECB mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := NewECBDecrypter(block)

	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	if string(plaintext[:33]) != "I'm back and I'm ringin' the bell" {
		t.Fatalf("ciphertext not properly decrypted, got", string(plaintext))
	}
}
