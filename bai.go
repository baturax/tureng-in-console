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
	req, err := http.NewRequest("GET", "https://tureng.com/tr/turkce-ingilizce/"+a, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0")
	req.Header.Set("Cookie", "VFRVREM%3d=dHI%3d; VFRESUNUSU9OQVJZ=ZW50cg%3d%3d; dm=1; cf_clearance=Xbhjc38NMEhPnRu8n5jhSc4XwRSOVIHtcwSJAU8HE7w-1757759919-1.2.1.1-uGIoyrayLxPKoJBFQmXHqceHs16DfPCRaQfQ3yW5ThouLsy5sLD_VRyISgX.W5fTs0BBuIgPtTLQYNQp1UCVgIx6jmuDZV2uYmEvDnKgk0I1zuA2wLYDZiRwlWg1VGfenTrmKXPaM1nihwfvsMPobZ8VtuBvsHmQ.zHZlny3GdA39nbgcMf4VTCJdrDsZvAAcSih_lSK.fraiTOfIN_2JJWf6Vrq9SyTWMfIk_LJ4vM; THI=bai=638933639217142460")

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
