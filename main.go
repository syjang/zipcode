package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fedesog/webdriver"
	"github.com/tealeg/xlsx"
	"github.com/tebeka/selenium"
)

func loadData(name string) {
	file, err := xlsx.OpenFile(name)
	if err != nil {
		log.Fatal(err)
	}

	targetSheet := file.Sheets[0]
	for _, row := range targetSheet.Rows {
		name := row.Cells[0].String()
		if name == "" {
			log.Println("종료!")
			break
		}
		log.Println(name, "검색중...")
		postnumber := findPostnumber(name)
		var cell *xlsx.Cell
		if len(row.Cells) > 2 {
			cell = row.Cells[1]
		}
		cell = row.AddCell()
		cell.SetValue(postnumber)
	}

	file.Save("output.xlsx")
}

func findPostnumber(name string) string {
	name = strings.Replace(name, " ", "%20", -1)
	address := "https://s.search.naver.com/n/csearch/content/eprender.nhn?where=nexearch&pkid=252&key=address_kor&q=" + name + "%20우편번호"
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return ""
	}
	// log.Println(address)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return ""
	}
	numberlist := []string{}
	doc.Find("td.tc").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		numberlist = append(numberlist, s.Find("strong").Text())
	})
	if len(numberlist) > 0 {
		return numberlist[0]
	}
	return ""
}

func loadDatav2(name string) {

	chromeDriver := webdriver.NewChromeDriver("./chromedriver")
	err := chromeDriver.Start()
	if err != nil {
		log.Println(err)
	}

	desired := webdriver.Capabilities{"Platform": "Windows"}
	required := webdriver.Capabilities{}
	session, err := chromeDriver.NewSession(desired, required)
	if err != nil {
		log.Println(err)
	}
	defer session.Delete()
	defer chromeDriver.Stop()

	file, err := xlsx.OpenFile(name)
	if err != nil {
		log.Fatal(err)
	}

	targetSheet := file.Sheets[0]
	for _, row := range targetSheet.Rows {
		name := row.Cells[0].String()
		if name == "" {
			log.Println("종료!")
			break
		}
		log.Println(name, "검색중...")
		postnumber := findPostNumberV2(name, session)
		var cell *xlsx.Cell
		if len(row.Cells) > 2 {
			cell = row.Cells[1]
		}
		cell = row.AddCell()
		cell.SetValue(postnumber)
	}

	file.Save("output.xlsx")
}

func findPostNumberV2(name string, session *webdriver.Session) string {
	url := "https://search.naver.com/search.naver?sm=tab_hty.top&where=nexearch&query=" + name + "%20우편번호&oquery=" + name + "%20우편번호&tqi=itQFowp0JXossfmUuT8ssssst8V-304727"
	err := session.Url(url)
	if err != nil {
		log.Println(err)
		return ""
	}

	span, err := session.FindElement(selenium.ByCSSSelector, "#loc-main-section-root > section > div > div.uOjIX > div > div.EzDG5 > div.g7j7K > div.Z4lWG > span:nth-child(1)")
	if err != nil {
		log.Println(err)
		return ""
	}
	str, err := span.Text()
	if err != nil {
		log.Println(err)
		return ""
	}
	reg, err := regexp.Compile("[0-9]+")
	if err != nil {
		log.Println(err)
		return ""
	}
	ret := reg.FindAllString(str, -1)

	return ret[0]
}

func main() {
	loadDatav2("input.xlsx")
}
