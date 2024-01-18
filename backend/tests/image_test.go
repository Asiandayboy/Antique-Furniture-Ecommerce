package tests

import (
	"backend/util"
	"testing"
)

func TestEncodeImageToBase64(t *testing.T) {
	data, _ := util.EncodeImageToBase64("../tests/test_images/tiger_maple1.jpg")

	if data == "" {
		t.Fatalf("Expected base64 encoded string, got empty string")
	}
}
