package main

import (
	"encoding/json"
	"fmt"
	"log"

	"regexp"
	"strings"
	"time"

	"context"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Sfnode struct {
	Created  time.Time  `json:"created"`
	Prefixes []Sfprefix `json:"prefixes" db:"prefixes"`
}

type Sfprefix struct {
	Region   string `json:"region" db:"region"`
	IPPrefix string `json:"ip_prefix" db:"ip_prefix"`
}

func scrapesf(url string) (html string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.OuterHTML(`#solutionWrapper > table > tbody > tr > td > div > div > table`, &html, chromedp.NodeVisible, chromedp.ByID),
	)

	if err != nil {
		log.Fatal("error!! ", err)
	}

	return
}

func parsehtml(html *string) (sfpref []Sfprefix) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}\/[0-9]{1,2}`)

	var reg string

	doc.Find("td").Each(func(_ int, s *goquery.Selection) {
		var sp Sfprefix

		val, ok := s.Attr("colspan")
		if ok && val == "2" {
			if s.Text() == "Australia (Public Cloud)" {
				reg = "Australia"
			} else if s.Text() == "Canada (Public Cloud)" {
				reg = "Canada"
			} else {
				reg = s.Text()
			}
		}
		match := re.FindString(s.Text())
		sp.IPPrefix = match
		sp.Region = reg

		if sp.IPPrefix != "" {
			sfpref = append(sfpref, sp)

		}
	})

	return
}

func main() {

	out := scrapesf("https://help.salesforce.com/articleView?id=000003652&type=1")

	//fmt.Println(out)
	prefixes := parsehtml(&out)

	x := Sfnode{
		Created:  time.Now(),
		Prefixes: prefixes,
	}

	outj, err := json.MarshalIndent(x, "", "    ")
	if err != nil {
		log.Fatal("JSON marshal error. ", err)
	}

	fmt.Println(string(outj))
}
