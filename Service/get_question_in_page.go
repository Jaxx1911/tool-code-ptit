package Service

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type ProblemInfo struct {
	problemName   string
	problemLink   string
	problemId     string
	problemStatus string
	problemTopic  string
	problemLevel  string
}

func GetQuestionInPage(cookie string, page int) ([]ProblemInfo, []ProblemInfo, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", apiURL+"/student/question"+"?page="+string(page), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	var completeQuestionList []ProblemInfo
	var incompleteQuestionList []ProblemInfo

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, nil, err
	}
	doc.Find(".status .ques__table__wrapper table.ques__table tbody").Each(func(i int, s *goquery.Selection) {
		problemName := strings.TrimSpace(s.Find("td:nth-child(4)").Text())
		problemLink, _ := s.Find("td:nth-child(4) a").Attr("href")
		problemID := strings.TrimSpace(s.Find("td:nth-child(3)").Text())
		problemTopic := strings.TrimSpace(s.Find("td:nth-child(6)").Text())
		problemDifficulty := strings.TrimSpace(s.Find("td:nth-child(7)").Text())

		problemInfo := ProblemInfo{
			problemName:  problemName,
			problemLink:  problemLink,
			problemId:    problemID,
			problemTopic: problemTopic,
			problemLevel: problemDifficulty,
		}

		if s.HasClass("bg--10th") {
			problemInfo.problemStatus = "Complete"
			completeQuestionList = append(completeQuestionList, problemInfo)
		} else {
			problemInfo.problemStatus = "Incomplete"
			incompleteQuestionList = append(incompleteQuestionList, problemInfo)
		}
	})
	return completeQuestionList, incompleteQuestionList, nil
}
