package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var text = tview.NewTextView()
var app = tview.NewApplication()

func main() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'Q' {
			app.Stop()
		}
		return event
	})

	text.SetDynamicColors(true).SetBorder(true).SetTitle("Press q to quit")

	if len(os.Args) < 2 {
		fmt.Println("Write a word:")
		reader := bufio.NewReader(os.Stdin)
		b, _ := reader.ReadString('\n')
		b = strings.TrimSpace(b)

		escapedInput := url.PathEscape(b)
		text.SetText(vet(escapedInput))
		app.SetRoot(text, true).Run()

	} else if os.Args[1] == "--help" || os.Args[1] == "-h" {
		help()

	} else {
		joinedArgs := strings.Join(os.Args[1:], " ")
		escapedInput := url.PathEscape(joinedArgs)
		text.SetText(vet(escapedInput))
		app.SetRoot(text, true).Run()
	}
}

func vet(a string) string {
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

	var results []string

	doc.Find("tr.tureng-manual-stripe-even, tr.tureng-manual-stripe-odd").Each(func(i int, s *goquery.Selection) {
		english := s.Find("td.en.tm a").Text()
		turkish := s.Find("td.tr.ts a").Text()

		if english != "" && turkish != "" {
			results = append(results, fmt.Sprintf("%d: %s -> %s\n", i+1, english, turkish))
		}
	})

	return strings.Join(results, "\n")
}

func help() {
	fmt.Println(`
-h, --help
	for help

Just write the word and over`)
}
