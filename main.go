package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)
type extractedJob struct{
	id string
	title string 
	jobCondition string 
	jobSector string	
	corpName string
}


var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=python"
// location := cleanString(card.Find(".area_job>.job_condition>span>a").Text()
	// 	title := cleanString(ca)
	func main() {
		var jobs []extractedJob
		c := make(chan []extractedJob)
		totalPages := getPages()
		
		for i := 0; i < totalPages; i++ {
			go getPage(i, c)
		}
		for i :=0; i<totalPages; i++{
			extractJobs :=<-c
			jobs = append(jobs, extractJobs...)
		}

		writeJobs(jobs)
		fmt.Println("Done,extracted", len(jobs))
	}
	

	// 각 페이지의 일자리를 추출하는 함수
	func getPage(page int,mainC chan<- []extractedJob) {
		var jobs []extractedJob
		c := make(chan extractedJob)
		// 필요한 url 만들기
		pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
		fmt.Println("Requesting", pageURL)
		res, err := http.Get(pageURL)
		checkErr(err)
		checkCode(res)
	
		defer res.Body.Close()
	
		doc, err := goquery.NewDocumentFromReader(res.Body)
		checkErr(err)
		item_recruit :=doc.Find(".item_recruit")
		// item_recruit :=doc.Find(".content")
			// fmt.Println(item_recruit)
		
		// 각 카드에서 일자리 정보 가져옴
		item_recruit.Each(func(i int, card *goquery.Selection) {
			// fmt.Println(card)
			go extractJob(card, c)
		})
		for i:=0; i < item_recruit.Length();i++{
			job := <-c
			jobs = append(jobs, job)
		}
		mainC <-jobs
	}
	func extractJob(card *goquery.Selection, c chan<- extractedJob){
		id, _ :=card.Attr("value");
		title := cleanString(card.Find(".job_tit>a").Text())
		jobCondition := cleanString(card.Find(".job_condition").Text())
		jobSector := cleanString(card.Find(".job_sector").Text())
		corpName := cleanString(card.Find(".corp_name").Text())
		// fmt.Println(id, title, jobCondition,jobSector,corpName)
		c<- extractedJob{
			id:   				id,
			title:     		title,
			jobCondition:	jobCondition,
			jobSector:		jobSector,
			corpName:			corpName,
		}
	}
	func cleanString(str string) string{
		// trimSpace 공백을 제거한다.
		// fields가 텍스트로만 이루어진 배열을 만든다
		// 
		return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
	}


	func getPages() int {
		pages := 0
		res, err := http.Get(baseURL)
		checkErr(err)
		checkCode(res)
	
		defer res.Body.Close()
	
		doc, err := goquery.NewDocumentFromReader(res.Body)
		checkErr(err)
	
		doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
			pages = s.Find("a").Length()
		})
	
		return pages
	}
	func writeJobs(jobs []extractedJob){
		file, err := os.Create("jobs.csv")
		checkErr(err)
		  
		w := csv.NewWriter(file)
		// writer를 생성하고, writer에 데이터를 입력하고, 모든 데이터를 파일에 저장
		defer w.Flush( )

		headers := []string{"Link","Title","jobCondition", "jobSector", "corpName"}
		 wErr := w.Write(headers)
		 checkErr(wErr)

		 for _, job := range jobs{
			jobSlice := []string{"https://www.saramin.co.kr/zf_user/jobs/relay/view?isMypage=no&rec_idx="+job.id, job.title, job.jobCondition, job.jobSector,job.corpName}
			jwErr := w.Write(jobSlice)
			checkErr(jwErr)
		 }
	}
	func checkErr(err error) {
		if err != nil {
			log.Fatalln(err)
		}
	}
	
	func checkCode(res *http.Response) {
		if res.StatusCode != 200 {
			log.Fatalln("Request failed with Status:", res.StatusCode)
		}
	}

