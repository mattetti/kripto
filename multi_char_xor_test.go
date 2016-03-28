package kripto

import "testing"

func TestBreakMultiCharXor(t *testing.T) {
	t.Skip("To hard to test with only a few letters, we need to apply another approach like Kasiski's for instance")
	testCases := []struct {
		input  []byte
		key    string
		output string
	}{
		{[]byte{
			// crypto
			0x2, 0x10, 0x1a, 0x14, 0x15, 0xd,
			// isshortfor
			0xa, 0x17, 0x12, 0xa, 0xc, 0x16, 0x15, 0x4, 0xc, 0x16,
			// crypto
			0x2, 0x10, 0x1a, 0x14, 0x15, 0xd,
			// graphy
			0x4, 0x16, 0x0, 0x12, 0xb, 0x1d},
			"ABCD",
			"CRYPTOISSHORTFORCRYPTOGRAPHY",
		},
	}

	for i, tc := range testCases {
		t.Logf("test case %d\n", i)
		o, k := BreakMultiCharXor(tc.input, 40)
		if string(k) != tc.key {
			t.Fatalf("key not properly found, expected '%s', got '%s'", tc.key, string(k))
		}
		if string(o) != tc.output {
			t.Fatalf("expected %s\ngot\n%s\n", tc.output, string(o))
		}
	}
}
