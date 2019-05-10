// https://www.devdungeon.com/content/web-scraping-go
// http://sandipbgt.com/2018/08/23/scraping-tutorial-with-golang/
// https://edmundmartin.com/scraping-google-with-golang/
// https://appdividend.com/2019/03/23/golang-receiver-function-tutorial-go-function-receivers-example/
// https://blog.heroku.com/neither-self-nor-this-receivers-in-go

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

type ZscalerRegion struct {
	ZRegion string `json:"region"`
	Element string `json:"-"`
}

type ZscalerNode struct {
	ZscalerRegion
	Prefix string `json:"ip_prefix"`
}

type Zscaler struct {
	Created  time.Time     `json:"created"`
	Prefixes []ZscalerNode `json:"prefixes"`
}

func getIPs(zr []*ZscalerRegion) ([]ZscalerNode, error) {

	response, err := http.Get("https://ips.zscaler.net/cenr")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	znode := make([]ZscalerNode, 0)
	for _, v := range zr {
		doc.Find(v.Element).Each(func(i int, s *goquery.Selection) {
			for _, line := range strings.Split(strings.TrimSuffix(s.Text(), "\n"), "\n") {
				var zn ZscalerNode
				zn.ZRegion = v.ZRegion
				zn.Prefix = strings.TrimSpace(line)
				znode = append(znode, zn)
			}
		})
	}

	return znode, nil
}

func main() {

	r := []*ZscalerRegion{
		{
			ZRegion: "Europe",
			Element: "#div_europe",
		},
		{
			ZRegion: "USCanada",
			Element: "#div_uscanada",
		},
		{
			ZRegion: "Asia",
			Element: "#div_asia",
		},
		{
			ZRegion: "Africa",
			Element: "#div_africa",
		},
		{
			ZRegion: "LatinAmerica",
			Element: "#div_latinamerica",
		},
	}

	z, _ := getIPs(r)

	x := new(Zscaler)
	cr := time.Now()
	x.Created = cr
	x.Prefixes = z

	out, err := json.Marshal(x)
	if err != nil {
		log.Fatal("JSON marshal error. ", err)
	}

	fmt.Println(string(out))

}
