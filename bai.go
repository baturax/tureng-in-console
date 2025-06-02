package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	q := strings.Join(os.Args[1:], " ")
	a := url.PathEscape(q)

	if len(os.Args) < 2 {
		help()
	} else if os.Args[1] == "--help" || os.Args[1] == "-h" {
		help()
	} else {
		translate(a)
	}

}

func translate(a string) {
	req, err := http.NewRequest("GET", "https://tureng.com/en/turkish-english/"+a, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("tr.tureng-manual-stripe-even, tr.tureng-manual-stripe-odd").Each(func(i int, s *goquery.Selection) {
		english := s.Find("td.en.tm a").Text()
		turkish := s.Find("td.tr.ts a").Text()

		if english != "" && turkish != "" {
			fmt.Printf("%d: %s -> %s\n", i+1, english, turkish)
		}
	})
}

func help() {
	fmt.Println(`
-h, --help
	for help

Just write the word and over`)
}
