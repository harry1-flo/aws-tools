package fs

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

var DIST_PATH = "dist"

type CSVParams struct {
	fs *os.File // multi file

	oneFS *os.File // one file
}

func NewCSV(name string) CSVParams {
	filename := getFilename(name)

	distDir, _ := os.Getwd()
	distDir = distDir + "../../../dist"

	fs, err := os.Create(filepath.Join(distDir, filename))
	if err != nil {
		panic("Failed to create CSV file: " + err.Error())
	}

	onefs, err := os.OpenFile(filepath.Join(distDir, "result.csv"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic("Failed to create CSV file: " + err.Error())
	}

	return CSVParams{
		fs:    fs,
		oneFS: onefs,
	}
}

func (c CSVParams) Write(data ...string) {
	c.fs.WriteString(strings.Join(data, ",") + "\n")
	c.fs.Sync() // 즉시 디스크에 쓰기
}

func (c CSVParams) OneFileWrite(data ...string) {
	c.oneFS.WriteString(strings.Join(data, ",") + "\n")
	c.oneFS.Sync() // 즉시 디스크에 쓰기
}

func (c CSVParams) End() {
	c.fs.Close()
	c.oneFS.Close()
}

func getFilename(name string) string {
	return time.Now().Format("20060102") + "_" + name + ".csv"
}
