package utils

import "os"

func ReadBytes(file *os.File, n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := file.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
