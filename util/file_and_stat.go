package util

import "os"

// FileAndStat opens a file and returns the file handle, its size, and any error encountered.
func FileAndStat(filename string) (*os.File, int64, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, 0, err
	}

	return file, fileInfo.Size(), nil
}
