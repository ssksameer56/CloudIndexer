package utils

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestFileReader(t *testing.T) {
	auth := "alice:pa55word"
	basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	fmt.Println(basicAuth)
}
