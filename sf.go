package main

import (
	"fmt"
	"log"
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

type Sfregion []string

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

func parsehtml(html string) (sfpref []Sfprefix) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		doc.Find("td").Each(func(j int, sel *goquery.Selection) {
			if sel.Text() == "RIPE" {
				fmt.Println(sel.Children().Length(), sel.Text(), sel.Next().Text())
			}
		})
	})

	return
}

func main() {
	out := scrapesf("https://help.salesforce.com/articleView?id=000003652&type=1")

	//fmt.Println(out)
	prefixes := parsehtml(out)
	fmt.Println(prefixes)

}
