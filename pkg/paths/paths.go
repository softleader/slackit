package paths

import "os"

// Exists checks the path exists
func Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// IsRegularFile checks the path is a regular file, not a directory
func IsRegularFile(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.Mode().IsRegular(), nil
}
