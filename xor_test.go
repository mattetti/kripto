package kripto

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

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
		{"6.txt", []int{2, 20, 7, 35, 4, 24, 28, 3, 32, 38}},
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
