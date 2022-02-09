package io

import (
	"bufio"
	"os"
)

// ReadLineByLine reads a file line by line and returns its contents and nil if successful.
// Returns nil and the error if one occurs.
func ReadLineByLine(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var content []string
	sc := bufio.NewScanner(file) // if lines length > 64k --> set sc.Buffer(buf, maxCapacity)
	for sc.Scan() {
		content = append(content, sc.Text())
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}
	return content, nil
}

// ReadNLinesFromFile reads the first n lines from a file and returns those and nil if successful.
// Returns nil and the error if one occurs.
func ReadNLinesFromFile(filename string, numLines int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var content []string
	sc := bufio.NewScanner(file) // if lines length > 64k --> set sc.Buffer(buf, maxCapacity)
	for i := 0; i < numLines; i++ {
		sc.Scan()
		content = append(content, sc.Text())
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}
	return content, nil
}
