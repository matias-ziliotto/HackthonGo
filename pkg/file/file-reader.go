package file

import (
	"bufio"
	"os"
	"strings"
)

func ReadFile(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var rowsScanned [][]string

	for scanner.Scan() {
		// Split text line into slice of strings and append to slice of rows
		rowsScanned = append(rowsScanned, strings.Split(scanner.Text(), "#$%#"))
	}

	return rowsScanned, nil
}
