package Service

import (
	"errors"
	"io/ioutil"
	"strings"
)

func GetSourceCode(problemCode string) (string, error) {
	allFile, err := ioutil.ReadDir("./SourceCode")
	if err != nil {
		return "", err
	}
	for _, file := range allFile {
		if strings.Contains(file.Name(), problemCode) {
			return file.Name(), nil
		}
	}
	return "", errors.New("Source code not found ")
}
