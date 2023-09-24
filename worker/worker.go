package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

type Result struct {
	Line    string
	LineNum int
	Path    string
}

type Results struct {
	Inner []Result
}

func NewResult(line string, lineNum int, path string) Result {
	return Result{line, lineNum, path}
}

func FindInFile(path, find string) *Results {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error openning file. err is :%s", err)
		return nil
	}
	defer file.Close()

	if isBinaryFile(file) {
		return nil
	}
	results := Results{Inner: make([]Result, 0)}

	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), find) {
			r := NewResult(scanner.Text(), lineNum, path)
			results.Inner = append(results.Inner, r)
			lineNum++
		}
	}
	if len(results.Inner) < 1 {
		return nil
	}
	return &results
}

func isBinaryFile(file *os.File) bool {

	mimeType, err := mimetype.DetectReader(file)
	if err != nil {
		fmt.Printf("could not detect %s mime type. err is : %s\n", file.Name(), err)
		return true
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Printf("could not seek to origin of %s. err is : %s\n ", file.Name(), err)
		return true
	}
	return mimeType.Is("application/x-mach-binary")
}
