package Service

import (
	"bytes"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Status string

const (
	A Status = "SUCCESS"
	B Status = "FAILED"
	C Status = "NOT_FOUND"
)

func SubmitCode(problemCode, cookie, csrf string) (Status, error) {
	const api = apiURL + "/student/solution"
	client := &http.Client{}
	var reqBody bytes.Buffer

	writer := multipart.NewWriter(&reqBody)
	sourceCodeName, err := GetSourceCode(problemCode)
	if err != nil {
		if err.Error() == "Source code not found " {
			return C, err
		} else {
			return B, err
		}
	}
	sourceCodePath := filepath.Join("./SourceCode/", sourceCodeName)
	sourceCode, err := os.Open(sourceCodePath)
	if err != nil {
		return B, err
	}
	defer sourceCode.Close()

	writer.WriteField("_token", csrf)
	writer.WriteField("question", problemCode)
	writer.WriteField("compiler", "2")

	part, err := writer.CreateFormFile("code_file", "code_file.cpp")
	if err != nil {
		log.Fatalf("Error creating form file: %v", err)
	}
	if _, err := io.Copy(part, sourceCode); err != nil {
		log.Fatalf("Error copying file content: %v", err)
	}
	writer.Close()
	req, err := http.NewRequest("POST", api, &reqBody)

	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Mobile Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return B, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return A, nil
	} else {
		return B, errors.New(strconv.Itoa(res.StatusCode))
	}
}
