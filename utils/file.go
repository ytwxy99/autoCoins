package utils

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var wLock *sync.RWMutex = new(sync.RWMutex)

// write []string into file
func WriteLines(lines []string, filePath string) error {
	wLock.Lock()
	f, err := os.Create(filePath)
	if err != nil {
		logrus.Error("create map file error: %v", err)
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range lines {
		fmt.Fprintln(w, fmt.Sprintf("%s", line))
	}

	err = w.Flush()
	wLock.Unlock()
	return err
}

// read lines from specified file
func ReadLines(filePath string) ([]string, error) {
	var lines []string

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}
