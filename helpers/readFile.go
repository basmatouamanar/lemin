package helpers

import (
	"os"
	"strings"
)

func ReadFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	content := string(data)
	s := strings.ReplaceAll(content, "\r", "")
	lines := strings.Split(s, "\n")
	return lines, nil
}
