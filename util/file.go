package util

import (
	"io"
	"io/ioutil"
	"mstats-new/logger"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ProgramPath() string {
	ex, err := os.Executable()

	if err != nil {
		logger.Error("Can't read program folder path")
		logger.Error(err.Error())
		return ""
	}

	exPath := filepath.Dir(ex)
	return exPath
}

func DownloadFile(path string, url string) error {
	out, err := os.Create(path)

	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return err
	}

	return nil
}

func ReadFileToInt(name string) int {
	bytes, err := ioutil.ReadFile(name)
	if err != nil {
		logger.Error("Something went wrong when reading file " + name)
		return 0
	}
	lines := strings.Split(string(bytes), "\n")
	num, err := strconv.Atoi(lines[0])

	if err != nil {
		logger.Error("Something went wrong when converting value in " + name)
		return 0
	}
	return num
}
