package slp

import (
	"encoding/base64"
	"os"
)

func FaviconFromFile(path string) (string, error) {
	// read the file
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()

	// read the file content
	buffer := make([]byte, fileSize)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	// convert the content to base64
	base64Str := base64.StdEncoding.EncodeToString(buffer)

	// return the data url
	return "data:image/png;base64," + base64Str, nil
}
