package utils

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

func InitLog(logPath string) {
	writer1 := &bytes.Buffer{}
	writer2 := os.Stdout
	writer3, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}

	logrus.SetOutput(io.MultiWriter(writer1, writer2, writer3))
}
