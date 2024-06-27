package Service

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

func GetCSRFCodePtit(cookie, problemCode string) (string, error) {
	api := apiURL + "/student/question/" + problemCode
	client := &http.Client{}

	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko)")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	csrf_token := doc.Find("input[name='_token']").AttrOr("value", "")
	return csrf_token, nil
}
