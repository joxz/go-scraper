package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Znode struct {
	Created  time.Time `json:"created"`
	Prefixes []Zprefix `json:"prefixes"`
}

type Zprefix struct {
	Region   string `json:"region"`
	Hostname string `json:"hostname"`
	Location string `json:"location"`
	IPPrefix string `json:"ip_prefix"`
}

type Zregion []string

func (z *Zregion) getregion() {
	response, err := http.Get("https://ips.zscaler.net/cenr")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		val, ok := s.Attr("id")
		if ok && s.HasClass("hidden") {
			*z = append(*z, strings.Replace(val, "div_", "", -1))
		}
	})
}

func main() {
	z := new(Zregion)
	z.getregion()
	fmt.Println(z)
}
