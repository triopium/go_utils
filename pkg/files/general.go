package files

import "os"

// IsFileReadable - check if file is readable. More reliable than using os.Stat.
func IsFileReadable(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	f.Close()
	return true, err
}
