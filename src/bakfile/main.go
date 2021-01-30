package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	_ = iota
	ERR_READ_CONFIG
	ERR_DST_FILE_IS_EMPTY
	ERR_OPEN_BAK_FILE

	ERR_LENGTH
)

var msgSlice []string
var config Config

func init() {
	msgSlice = make([]string, ERR_LENGTH)
	msgSlice[ERR_READ_CONFIG] = "read config failed."
	msgSlice[ERR_DST_FILE_IS_EMPTY] = "target file is abnormal. "
	msgSlice[ERR_OPEN_BAK_FILE] = "create bak file failed."

}

func main() {
	if !readConfig() {
		exit(ERR_READ_CONFIG)
	}

	srcFileBytes := readFile(config.File)
	if len(srcFileBytes) == 0 {
		exit(ERR_DST_FILE_IS_EMPTY)
	}

	bakFilePath := getBakFilePath(config.File)
	bakFileBytes := readFile(bakFilePath)

	if bytes.Equal(srcFileBytes, bakFileBytes) {
		return
	}

	fd, err := os.Create(bakFilePath)
	if err != nil {
		fmt.Println(err)
		exit(ERR_OPEN_BAK_FILE)
	}

	_, err = fd.Write(srcFileBytes)
	if err != nil {
		fmt.Println(err)
		exit(ERR_OPEN_BAK_FILE)
	}

	fd.Close()
}

type Config struct {
	File string `json:"file"`
}

// config.json
// {"file":"PathToFile"}
func readConfig() bool {
	b := readFile("config.json")
	if b == nil {
		fmt.Println("read config.json error.")
		return false
	}

	err := json.Unmarshal(b, &config)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func readFile(filePath string) []byte {
	fd, err := os.Open(filePath)
	if err != nil {
		return nil
	}

	var b []byte
	const buffSize = 10

	for {
		buf := make([]byte, buffSize)

		n, err := fd.Read(buf)
		if err == io.EOF {
			break
		}

		if n < buffSize {
			b = append(b, buf[:n]...)
			break
		} else {
			b = append(b, buf...)
		}

	}

	return b
}

func exit(codeEnum int) {
	msg := "error"

	if codeEnum < ERR_LENGTH {
		if msgSlice[codeEnum] != "" {
			msg = msgSlice[codeEnum]
		}
	}

	fmt.Println(msg)
	os.Exit(codeEnum)
}

func getBakFilePath(srcFilePath string) string {
	var extStartPos int
	extName := filepath.Ext(srcFilePath)
	srcFileBytes := []byte(srcFilePath)

	if extName == "" {
		extStartPos = len(srcFileBytes)
	} else {
		extStartPos = strings.Index(srcFilePath, extName)
	}

	return string(append(srcFileBytes[:extStartPos], []byte(".bak"+extName)...))
}
