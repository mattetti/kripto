package kripto

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestBreakSingleCharXor(t *testing.T) {

	testCases := []struct {
		input  string
		output string
		fn     func([]byte) []byte
		key    byte
	}{
		{
			"1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736",
			"Cooking MC's like a pound of bacon",
			DeHex,
			'X',
		},
		{
			"7b5a4215415d544115415d5015455447414c155c46155f4058455c5b523f",
			"Now that the party is jumping\n",
			DeHex,
			'5',
		},
		// Base64 test
		{
			"aVhTWl5FCkNZCkxfRAYKQ1lEDV4KQ14V",
			"Crypto is fun, isn't it?",
			DeBase64,
			42,
		},
		// no processing test
		{
			string([]byte{0xe, 0x3f, 0x34, 0x3d, 0x39, 0x22, 0x6d, 0x24, 0x3e, 0x6d, 0x2b, 0x38, 0x23, 0x61, 0x6d, 0x24, 0x3e, 0x23, 0x6a, 0x39, 0x6d, 0x24, 0x39, 0x72}),
			"Crypto is fun, isn't it?",
			nil,
			'M',
		},
	}

	scorer := &EnglishScorer{WithSpace: true}
	for i, tc := range testCases {
		t.Logf("test case %d\n", i)
		if o, k := BreakSingleCharXor([]byte(tc.input), tc.fn, scorer); string(o) != tc.output || k != tc.key {
			t.Fatalf("expected to get:\b'%s' using key %s\nbut got\n'%s'\n from key %s\n", tc.output, string(tc.key), o, string(k))
		}
	}
}

func TestMultiCharXor(t *testing.T) {
	testCases := []struct {
		input  string
		key    string
		output []byte
	}{
		{"CRYPTOISSHORTFORCRYPTOGRAPHY",
			"ABCD",
			// crypto
			[]byte{
				// crypto
				0x2, 0x10, 0x1a, 0x14, 0x15, 0xd,
				// isshortfor
				0xa, 0x17, 0x12, 0xa, 0xc, 0x16, 0x15, 0x4, 0xc, 0x16,
				// crypto
				0x2, 0x10, 0x1a, 0x14, 0x15, 0xd,
				// graphy
				0x4, 0x16, 0x0, 0x12, 0xb, 0x1d},
		},
	}

	for i, tc := range testCases {
		t.Logf("test case %d\n", i)
		o := MultiCharXor([]byte(tc.input), []byte(tc.key))
		if bytes.Compare(o, []byte(tc.output)) != 0 {
			t.Fatalf("expected %#v\ngot\n%#v\n", []byte(tc.output), o)
		}
		o = MultiCharXor(o, []byte(tc.key))
		if bytes.Compare(o, []byte(tc.input)) != 0 {
			t.Fatalf("expected to go back to %#v\ngot\n%#v\n", []byte(tc.input), o)
		}
	}
}

func TestGuessMultiCharXorKeySize(t *testing.T) {
	testCases := []struct {
		file  string
		sizes []int
	}{
		{"6.txt", []int{3, 8, 32, 28, 38, 25, 35, 20, 31, 24}},
	}

	for i, tc := range testCases {
		t.Logf("test case %d\n", i)
		data, err := ioutil.ReadFile(fixturePath(tc.file))
		if err != nil {
			t.Fatal(err)
		}
		if o := GuessMultiCharXorKeySize(data, 40); !reflect.DeepEqual(o, tc.sizes) {
			t.Fatalf("expected %v\ngot\n%v\n", tc.sizes, o)
		}
	}
}

func TestMostLikelyXorKey(t *testing.T) {
	t.Skip("To hard to test with only a few letters, we need to apply Kasiski approach")
	testCases := []struct {
		input []byte
		k     byte
	}{
		0: {[]byte{0x2, 0x15, 0x12, 0x15, 0x2, 0x15, 0x0}, 'A'},
		1: {[]byte{0x10, 0xd, 0xa, 0x4, 0x10, 0xd, 0x12}, 'B'},
		2: {[]byte{0x1a, 0xa, 0xc, 0xc, 0x1a, 0x4, 0xb}, 'C'},
		3: {[]byte{0x14, 0x17, 0x16, 0x16, 0x14, 0x16, 0x1d}, 'D'},
	}

	for i, tc := range testCases {
		t.Logf("test case %d\n", i)
		if k := MostLikelyXorKey(tc.input); k != tc.k {
			t.Fatalf("expected key %s\ngot\n%s\n", string(tc.k), string(k))
		}
	}
}
