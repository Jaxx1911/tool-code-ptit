package Service

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type MySession struct {
	Client  *http.Client
	Host    string
	Headers map[string]string
}

func NewMySession(host string) *MySession {
	headers := map[string]string{
		"Accept":             "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		"Content-Type":       "application/x-www-form-urlencoded", // Đặt Content-Type thành application/x-www-form-urlencoded
		"Referer":            host + "/",
		"Sec-Ch-Ua":          `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`,
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": `"Windows"`,
	}

	return &MySession{
		Client:  &http.Client{},
		Host:    host,
		Headers: headers,
	}
}

func (s *MySession) Get(endpoint string) error {
	url := s.Host + endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for key, value := range s.Headers {
		req.Header.Set(key, value)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (s *MySession) Post(endpoint string, data url.Values) error {
	url := s.Host + endpoint

	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	for key, value := range s.Headers {
		req.Header.Set(key, value)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func structToURLValues(data interface{}) (url.Values, error) {
	values := url.Values{}
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typeField := t.Field(i)
		tag := typeField.Tag.Get("form")

		if tag != "" {
			values.Set(tag, fmt.Sprintf("%v", field.Interface()))
		}
	}
	return values, nil
}
