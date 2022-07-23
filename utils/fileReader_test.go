package utils

import (
	"io"
	"testing"
)

func TestFileReader(t *testing.T) {
	ReadFile(&io.PipeReader{})
}
