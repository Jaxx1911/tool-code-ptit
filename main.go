package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"tool-crawl-code-ptit/Service"
)

func delay() {
	randomDelay := rand.Intn(delayHigh-delayLow+1) + delayLow
	fmt.Printf("Waiting %d minutes to submit next code!\n", randomDelay)
	time.Sleep(time.Duration(randomDelay) * time.Minute)
}

const (
	delayLow  = 1
	delayHigh = 5
	limitTask = 10
)

var (
	countTask = 0
)

func main() {
	cookie, _, err := Service.GetCodePtitCookie()
	if err != nil {
		return
	}
	var completeProblemList, incompleteProblemList []Service.ProblemInfo

	// lọc lại câu hỏi
	for i := 1; i <= 4; i++ {
		completeProblem, incompleteProblem, err := Service.GetQuestionInPage(cookie, i)
		if err != nil {
			return
		}
		completeProblemList = append(completeProblemList, completeProblem...)
		incompleteProblemList = append(incompleteProblemList, incompleteProblem...)
	}

	problemMap := make(map[string]Service.ProblemInfo)
	for _, problem := range append(incompleteProblemList, completeProblemList...) {
		problemMap[problem.ProblemId] = problem
	}
	var problemArray []Service.ProblemInfo
	for _, problem := range problemMap {
		problemArray = append(problemArray, problem)
	}
	incompleteProblemList = nil
	completeProblemList = nil
	for _, problem := range problemArray {
		if problem.ProblemStatus == "Incomplete" {
			incompleteProblemList = append(incompleteProblemList, problem)
		} else if problem.ProblemStatus == "Complete" {
			completeProblemList = append(completeProblemList, problem)
		}
	}

	//
	randomIndexCsrf := rand.Intn(len(incompleteProblemList))
	randomProblemIDCsrf := incompleteProblemList[randomIndexCsrf].ProblemId
	submitCsrf, err := Service.GetCSRFCodePtit(randomProblemIDCsrf, cookie)
	if err != nil {
		log.Fatalf("Error getting CSRF token: %v", err)
	}

	fmt.Printf("\x1b[32m%s\x1b[0m", "Login success!\n")
	fmt.Printf("> Total problems: %d\n", len(problemArray))
	fmt.Printf("> Incomplete problems: %d\n", len(incompleteProblemList))
	fmt.Printf("> Complete problems: %d\n", len(completeProblemList))

	delay()

	fmt.Printf("\x1b[34m%s\x1b[0m", fmt.Sprintf("Start auto submit code! Delay: %d - %d minutes\n", delayLow, delayHigh))

	for countTask < limitTask {
		randomIndex := rand.Intn(len(incompleteProblemList))
		randomProblemID := incompleteProblemList[randomIndex].ProblemId
		randomProblemName := incompleteProblemList[randomIndex].ProblemName

		fmt.Printf("\x1b[33m%s\x1b[0m", fmt.Sprintf("Submitting problemID: %s, problemName: %s\n", randomProblemID, randomProblemName))

		result, err := Service.SubmitCode(randomProblemID, cookie, submitCsrf)
		if err != nil {
			return
		}
		if result == Service.A {
			countTask++
			fmt.Printf("\x1b[33m%s\x1b[0m", fmt.Sprintf("Submitted %d/%d times\n", countTask, limitTask))
		}

		if result != Service.C {
			randomDelay := rand.Intn(delayHigh-delayLow+1) + delayLow
			fmt.Printf("\x1b[33m%s\x1b[0m", fmt.Sprintf("Waiting %d minutes to submit next code!\n", randomDelay))
			delay()
		}
	}

	fmt.Printf("\x1b[32m%s\x1b[0m", "Đã làm đủ bài rồi! Mai chạy tiếp !\n")
}
