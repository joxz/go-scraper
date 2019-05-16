package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Znode struct {
	Created  time.Time `json:"created"`
	Prefixes []Zprefix `json:"prefixes" db:"prefixes"`
}

type Zprefix struct {
	Region   string `json:"region" db:"region"`
	Hostname string `json:"hostname" db:"hostname"`
	Location string `json:"location" db:"location"`
	IPPrefix string `json:"ip_prefix" db:"ip_prefix"`
}

type Zregion []string

func getselection(url string) (*goquery.Document, error) {
	response, err := http.Get(url)

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func getregion(doc *goquery.Document) (z *Zregion) {

	z = new(Zregion)

	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		val, ok := s.Attr("id")
		if ok && s.HasClass("hidden") {
			*z = append(*z, strings.Replace(val, "div_", "", -1))
		}
	})
	return
}

func getprefix(z *Zregion, doc *goquery.Document) (zpref []Zprefix) {

	for _, v := range *z {
		div := fmt.Sprintf("#div_%s", v)
		doc.Find(div).Each(func(i int, s *goquery.Selection) {
			for _, line := range strings.Split(strings.TrimSuffix(s.Text(), "\n"), "\n") {
				var zp Zprefix
				zp.Region = v
				zp.IPPrefix = strings.TrimSpace(line)

				doc.Find("#block-system-main > div > div > div > article > table > tbody > tr > td").Each(func(j int, sel *goquery.Selection) {
					if strings.Contains(sel.Text(), string(line)) {
						zp.Hostname = strings.Replace(sel.Next().Text(), "\t", "", -1)
						zp.Location = sel.Prev().Text()
					}
				})
				zpref = append(zpref, zp)
			}

		})
	}
	return
}

func main() {
	selection, err := getselection("https://ips.zscaler.net/cenr")
	regions := getregion(selection)

	x := Znode{
		Created:  time.Now(),
		Prefixes: getprefix(regions, selection),
	}

	out, err := json.MarshalIndent(x, "", "    ")
	if err != nil {
		log.Fatal("JSON marshal error. ", err)
	}

	fmt.Println(string(out))
}
