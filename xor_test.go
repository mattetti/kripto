package kripto

import (
	"io/ioutil"
	"reflect"
	"testing"
)

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
