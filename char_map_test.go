package kripto

import "testing"

func TestLooksEnglish(t *testing.T) {
	Debug = false

	testCases := []struct {
		input     []string
		winnerIDX int
	}{
		{[]string{
			`q]]Y[\UqA^[YWSB]G\V]TPSQ]\`,
			`Cooking MC's like a pound of bacon`,
			`Q}}y{|u2_Q5a2~{yw2s2b}g|v2}t2psq}|`,
			"Dhhlni`'JD t'knlb'f'whric'ha'efdhi",
		},
			1},
	}

	for i, tc := range testCases {
		t.Logf("test case %d\n", i)
		maps := make(CharMaps, len(tc.input))
		for j, s := range tc.input {
			maps[j] = NewCharMap([]byte(s))
		}
		winner := maps.MostEnglish()
		if winner != maps[tc.winnerIDX] {
			t.Fatalf("expected %+s\ngot\n%+s\n", maps[tc.winnerIDX], winner)
		}
	}
}
