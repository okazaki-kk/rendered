package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

var chartFileSeparator = []byte("---\n# Source: ")
var outputDir string

func init() {
	flag.StringVar(&outputDir, "output-dir", "output", "Output directory")
}

// scanChartFile is a split function for a Scanner that returns each chart file from helm template output
func scanChartFile(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Find the next chart file separator
	if i := bytes.Index(data, chartFileSeparator); i >= 0 {
		return i + len(chartFileSeparator), data[0:i], nil
	}

	// If we're at the end of the input, return the remaining data
	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

// splitContent splits the content into the first line and the rest of the content
// The first line is the file name and the rest of the content is the content of the chart file
func splitContent(content string) (string, string) {
	if i := strings.Index(content, "\n"); i >= 0 {
		return content[:i], content[i+1:]
	}
	return "", ""
}

func main() {
	flag.Parse()

	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		log.Printf("Output directory %s already exists", outputDir)
	}

	os.Mkdir(outputDir, 0755)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(scanChartFile)

	// Skip the first chart file
	scanner.Scan()

	for scanner.Scan() {
		content := scanner.Text()
		fileName, fileContent := splitContent(content)
		fmt.Println("filename: ", fileName)

		outputFilePath := path.Join(outputDir, fileName)
		outputFileDir := path.Dir(outputFilePath)
		if _, err := os.Stat(outputFileDir); os.IsNotExist(err) {
			os.MkdirAll(outputFileDir, 0755)
		}

		if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
			log.Printf("Writing %s\n", outputFilePath)
			file, err := os.Create(outputFilePath)
			if err != nil {
				log.Fatalf("Failed to create file %s: %v", outputFilePath, err)
			}
			defer file.Close()
			file.WriteString(fileContent)
		} else {
			log.Printf("File %s already exists, overwrrite", outputFilePath)
			file, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				log.Fatalf("Failed to create file %s: %v", outputFilePath, err)
			}
			defer file.Close()
			file.WriteString(fileContent)
		}
		fmt.Printf("Writing %s\n", outputFilePath)
	}
}
