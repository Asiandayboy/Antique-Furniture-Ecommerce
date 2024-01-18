package util

import (
	"encoding/base64"
	"os"
)

func EncodeImageToBase64(filepath string) (string, error) {
	// convert image file to binary data
	binData, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	// convert binary data to base64 encoded string
	base64Encoded := base64.StdEncoding.EncodeToString(binData)

	return base64Encoded, nil
}
