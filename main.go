package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	getPages()
}

func getPages() int {
	// res, err := http.Get(baseURL)
	req, rErr := http.NewRequest("GET", baseURL, nil)
	checkErr(rErr)
	purl, err := url.Parse(baseURL)
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(purl)}}
	res, err := client.Do(req)
	checkErr(err)
	checkCode(res)




	// 메모리가 새는 것을 막는다
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination")

	return 0
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

