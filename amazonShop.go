package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tealeg/xlsx"
)

type item struct {
	name          string
	url           string
	buyboxgone    bool
	lightningdeal bool
	otherseller   [][]string
	rank          []string
}

//var findSubURL regexp.MustCompile("gp/offer-listing/([0-9a-zA-z?=;/])/")
/*func ExampleScrape() {
	// Request the HTML page.
	res, err := http.Get("https://www.amazon.com/PicassoTiles-Inflatable-Bouncing-Playhouse-Basketball/dp/B071W2H64Z?language=en_US")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(doc)
	// Find the review items
	doc.Find(".a-color-secondary a-size-base prodDetSectionEntry").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("span").Text()
		title := s.Find(".priceblock_ourprice").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
}*/

func queryhtmlToResp(url string) *http.Response {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36")
	req.Header.Set("cookie", "session-id=137-3436788-4375040")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return nil
	}
	/*defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)*/

	return resp
}

func queryhtmlToString(url string) (*http.Response, string) {
	/*resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http get error", err)
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error", err)
		return string(body)
	}*/
	resp := queryhtmlToResp(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read all", err)
	}
	return resp, string(body)
}

func otherSeller(url string) [][]string {
	var result [][]string
	findotherSeller, _ := regexp.Compile("from seller\\s*([a-zA-Z0-9 $.-]+)")
	//findotherPrice, _ := regexp.Compile("and price $\\s*([0-9.]+)<")
	url = "https://www.amazon.com/gp/offer-listing/" + url
	_, response := queryhtmlToString(url)
	//ioutil.WriteFile("otherSeller", []byte(response), 0777)
	for _, value := range findotherSeller.FindAllStringSubmatch(response, -1) {
		if strings.Split(value[1], " and price ")[0] != "Amazon Warehouse" &&
			strings.Split(value[1], " and price ")[0] != "Primo Super-store" &&
			strings.Split(value[1], " and price ")[0] != "KickBOT" &&
			strings.Split(value[1], " and price ")[0] != "SPORTSBOT" {
			result = append(result, strings.Split(value[1], " and price "))
		}
	}
	//fmt.Println(result)
	return result
}

func buyBoxGone(html string) bool {
	if strings.Contains(html, "Currently Unavailable") {
		return true
	}
	return false
}
func lightningDeal(html string) bool {
	if strings.Contains(html, "Lightning Deal") ||
		strings.Contains(html, "Lightning deal") ||
		strings.Contains(html, "lightning deal") {
		return true
	}
	return false
}

func test(resp *http.Response) [][]string {
	var result [][]string
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(doc)
	// Find the review items
	doc.Find("a-color-secondary a-size-base prodDetSectionEntry").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("span").Text()
		title := s.Find(".priceblock_ourprice").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
	return result
}

func bestSellersRank(html string) []string {
	var result []string
	rank := regexp.MustCompile("#\\s*([0-9,]+)")                                           //find rank
	rankin1 := regexp.MustCompile("in\\s*([0-9A-Za-z &-]+)\\s*<<a href='/gp/bestsellers/") //find rank in
	rankin2 := regexp.MustCompile("'>\\s*([0-9A-Za-z &'-]+)</a></span>")                   //find rank in
	for _, value := range strings.Split(html, "\n") {
		if strings.Contains(value, "<a href='/gp/bestsellers/") {
			value = strings.TrimSpace(strings.ReplaceAll(value, "(", "<"))
			if len(rankin1.FindStringSubmatch(value)) > 1 {
				if len(rank.FindStringSubmatch(value)) > 1 {
					result = append(result, "#"+rank.FindStringSubmatch(value)[1]+" in "+rankin1.FindStringSubmatch(value)[1])
				}
			} else if len(rankin2.FindStringSubmatch(value)) > 1 {
				if len(rank.FindStringSubmatch(value)) > 1 {
					result = append(result, "#"+rank.FindStringSubmatch(value)[1]+" in "+rankin2.FindStringSubmatch(value)[1])
				}
			} else {
				result = append(result, value)
			}
		}
	}
	return result
}
func output(output []item) {
	file, err := xlsx.OpenFile("format.xlsx")
	if err != nil {
		panic(err)
	}
	first := file.Sheets[0]
	row := first.AddRow()
	row.SetHeightCM(1)
	for _, value := range output {
		cell := row.AddCell()
		cell.Value = value.name
		for _, j := range value.rank {
			cell = row.AddCell()
			cell.Value = j
		}
		for _, j := range value.otherseller {
			cell = row.AddCell()
			cell.Value = strings.Join(j, " for ")
		}
		cell = row.AddCell()
		if value.buyboxgone {
			cell.Value = "true"
		} else {
			cell.Value = "false"
		}
		cell = row.AddCell()
		if value.lightningdeal {
			cell.Value = "true"
		} else {
			cell.Value = "false"
		}
		row = first.AddRow()
		row.SetHeightCM(1)
	}

	err = file.Save("Output.xlsx")
	if err != nil {
		panic(err)
	}
}

func webCrawler() {
	findSubURL, _ := regexp.Compile("offer-listing/\\s*([0-9a-zA-z?=;/_&]+)")

	content, err := ioutil.ReadFile("url")
	if err != nil {
		fmt.Println("Read file error: ", err)
	}
	body := string(content)
	parseurl := strings.Split(body, "\n")

	var product []item
	for _, value := range parseurl {
		var temp item
		temp.name = strings.TrimSpace(strings.Split(value, "\\")[0])
		temp.url = strings.TrimSpace(strings.Split(value, "\\")[1]) + "?language=en_US"
		fmt.Println("Collecting url", temp.url)
		_, html := queryhtmlToString(temp.url)
		//fmt.Println(html)
		if len(findSubURL.FindStringSubmatch(html)) > 1 {
			temp.otherseller = otherSeller(findSubURL.FindStringSubmatch(html)[1])
		}
		temp.buyboxgone = buyBoxGone(html)
		temp.lightningdeal = lightningDeal(html)
		temp.rank = bestSellersRank(html)
		product = append(product, temp)
	}
	//fmt.Println(product)
	output(product)
}

func main() {
	webCrawler()
}
