package main

import (
	"bufio"
	"bytes"
	"flag"
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

	// Create output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(scanChartFile)

	// Skip the first chart file
	scanner.Scan()

	for scanner.Scan() {
		content := scanner.Text()
		fileName, fileContent := splitContent(content)

		outputFilePath := path.Join(outputDir, fileName)
		outputFileDir := path.Dir(outputFilePath)
		if _, err := os.Stat(outputFileDir); os.IsNotExist(err) {
			os.MkdirAll(outputFileDir, 0755)
		}

		var file *os.File
		if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
			if file, err = os.Create(outputFilePath); err != nil {
				log.Fatalf("Failed to create file %s: %v", outputFilePath, err)
			}
		} else {
			if file, err = os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755); err != nil {
				log.Fatalf("Failed to open file %s: %v", outputFilePath, err)
			}
		}

		defer file.Close()
		if _, err := file.WriteString(fileContent); err != nil {
			log.Fatalf("Failed to write to file %s: %v", outputFilePath, err)
		}
	}
	log.Printf("Rendered helm chart files to %s\n", outputDir)
}
