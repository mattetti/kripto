package kripto

import "testing"

func TestHammingDiff(t *testing.T) {
	testCases := []struct {
		a    string
		b    string
		diff int
	}{
		{"this is a test", "wokka wokka!!!", 37},
	}

	for i, tc := range testCases {
		t.Logf("test case %d\n", i)
		if o := HammingDiff([]byte(tc.a), []byte(tc.b)); o != tc.diff {
			t.Fatalf("expected %d\ngot\n%d\n", tc.diff, o)
		}
	}
}
