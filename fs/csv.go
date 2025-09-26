package fs

import (
	"os"
	"strings"
	"time"
)

var DIST_PATH = "dist"

type CSVParams struct {
	fs *os.File
}

func NewCSV(name string) CSVParams {
	filename := getFilename(name)
	fs, err := os.Create(DIST_PATH + "/" + filename)
	if err != nil {
		panic("Failed to create CSV file: " + err.Error())
	}

	return CSVParams{
		fs: fs,
	}
}

func (c CSVParams) Write(data ...string) {
	c.fs.WriteString(strings.Join(data, ",") + "\n")
	c.fs.Sync() // 즉시 디스크에 쓰기
}

func (c CSVParams) End() {
	c.fs.Close()
}

func getFilename(name string) string {
	return time.Now().Format("20060102") + "_" + name + ".csv"
}
