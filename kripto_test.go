package kripto

import "path/filepath"

func fixturePath(filename string) string {
	fixturePath, _ := filepath.Abs("./fixtures")
	path := filepath.Join(fixturePath, filename)
	return path
}
