package hash

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReplaceHashes(rawInputFile string) {
	hashMap := make(map[string]string)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			hash := parts[0]
			nameString := parts[1]
			hashMap[nameString] = hash
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading hashes from stdin: %v\n", err)
		return
	}

	inputFile, err := os.Open(rawInputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer inputFile.Close()

	tempFile, err := os.CreateTemp("", "temp_manifest_*.json")
	if err != nil {
		fmt.Printf("Error creating temporary file: %v\n", err)
		return
	}
	defer os.Remove(tempFile.Name())

	scanner = bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(tempFile)

	for scanner.Scan() {
		line := scanner.Text()
		for nameString, hash := range hashMap {
			placeholder := nameString + "_shasum"
			if strings.Contains(line, placeholder) {
				line = strings.ReplaceAll(line, placeholder, hash)
			}
		}
		writer.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	writer.Flush()
	tempFile.Close()
	inputFile.Close()

	if err := os.Rename(tempFile.Name(), rawInputFile); err != nil {
		fmt.Printf("Error replacing original file: %v\n", err)
		return
	}
}
