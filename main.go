package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tealeg/xlsx"
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

func main() {
	loadData("input.xlsx")
}
