package Service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
)

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"_token"`
}

var host = "https://code.ptit.edu.vn/login"

func GetAccount() (Account, error) {
	var acc Account
	content, err := ioutil.ReadFile("./Data/account.txt")
	if err != nil {
		return acc, err
	}
	er := json.Unmarshal(content, &acc)
	if er != nil {
		return acc, er
	}
	return acc, nil
}

func GetCookieAndToken() ([](*http.Cookie), string, error) {
	client := &http.Client{}

	// Tạo một yêu cầu GET
	req, err := http.NewRequest("GET", "https://code.ptit.edu.vn/login", nil)
	if err != nil {
		return nil, "", err
	}

	// Gửi yêu cầu
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	tokenValue := doc.Find("input[name='_token']").AttrOr("value", "")

	// Lấy cookie từ phản hồi
	cookies := resp.Cookies()

	return cookies, tokenValue, nil
}

func LoginService() error {
	client := &http.Client{}
	account, err := GetAccount()
	if err != nil {
		return err
	}
	cookie, tk, err := GetCookieAndToken()
	if err != nil {
		return err
	}
	account.Token = tk

	body, err := json.Marshal(account)
	if err != nil {
		return err
	}

	res, err := http.NewRequest("POST", host, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	res.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Mobile Safari/537.36")

	for _, ck := range cookie {
		res.AddCookie(ck)
	}
	resp, err := client.Do(res)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBod, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(respBod))
	return nil
}
