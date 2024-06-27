package Service

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"_token"`
}

var apiURL = "https://code.ptit.edu.vn"
var userAgent = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Mobile Safari/537.36"

func GetAccount(token string) (Account, error) {
	acc := Account{
		Token: strings.ReplaceAll(token, " ", ""),
	}
	content, err := ioutil.ReadFile("./Data/account.txt")
	if err != nil {
		return acc, err
	}
	er := json.Unmarshal(content, &acc)
	if er != nil {
		return acc, er
	}
	fmt.Println(acc)
	return acc, nil

}

func GetCodePtitCookie() (string, string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", apiURL+"/login", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	var csrf_set_cookie string
	if len(res.Header["Set-Cookie"]) > 1 {
		csrf_set_cookie = res.Header["Set-Cookie"][1]
	} else {
		csrf_set_cookie = ""
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	tokenValue := doc.Find("input[name='_token']").AttrOr("value", "")
	account, err := GetAccount(tokenValue)
	if err != nil {
		return "", "", err
	}

	loginReq, err := http.NewRequest("POST", apiURL+"/login", strings.NewReader(fmt.Sprintf("_token=%s&username=%s&password=%s", account.Token, account.Username, account.Password)))
	if err != nil {
		return "", "", err
	}
	loginReq.Header.Set("User-Agent", userAgent)
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginReq.Header.Set("Cookie", csrf_set_cookie)

	loginRes, err := client.Do(loginReq)
	if err != nil {
		return "", "", err
	}
	defer loginRes.Body.Close()
	cookie := loginRes.Header["Set-Cookie"][1]
	fmt.Println("Cookie cuoi: " + cookie)
	fmt.Println("CSRF:" + csrf_set_cookie)
	return cookie, csrf_set_cookie, nil
}
