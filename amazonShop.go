package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tealeg/xlsx"
)

type item struct {
	name                 string
	url                  string
	price                string
	buyboxgone           bool
	lightningdeal        bool
	otherseller          [][]string
	rank                 []string
	leaderboard          []string
	currentlyunavailable bool
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
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
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

func otherSeller(url string) ([][]string, string) { //return seller who is not Amazon,Primo Super-store,KickBOT and SPORTSBOT and the cheapest price.
	var result [][]string
	//var cheapest string
	var cheap float64
	findotherSeller, _ := regexp.Compile("from seller\\s*([a-zA-Z0-9 $.-]+)")
	//findotherPrice, _ := regexp.Compile("and price $\\s*([0-9.]+) <")
	url = "https://www.amazon.com/gp/offer-listing/" + url
	//fmt.Println(url)
	_, response := queryhtmlToString(url)
	//ioutil.WriteFile("otherSeller", []byte(response), 0777)
	for _, value := range findotherSeller.FindAllStringSubmatch(response, -1) {
		if strings.Split(value[1], " and price ")[0] != "Amazon Warehouse" &&
			strings.Split(value[1], " and price ")[0] != "Primo Super-store" &&
			strings.Split(value[1], " and price ")[0] != "KickBOT" &&
			strings.Split(value[1], " and price ")[0] != "SPORTSBOT" {
			result = append(result, strings.Split(value[1], " and price "))
		} else if strings.Split(value[1], " and price ")[0] != "Amazon Warehouse" {
			if len(strings.Split(value[1], " and price $")) > 1 {
				s, err := strconv.ParseFloat(strings.Split(value[1], " and price $")[1], 32)
				if err != nil {
					fmt.Println("Price parseing error ", err)
				}
				s, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", s), 64)
				if s < cheap || cheap == 0 {
					cheap = s
				}
			}
		}
	}

	//fmt.Println(cheap)
	return result, fmt.Sprintf("%.2f", cheap)
}
func currentlyUnavailable(html string) bool {
	if strings.Contains(html, "Currently Unavailable") || strings.Contains(html, "Currently unavailable") {
		return true
	}
	return false
}
func buyBoxGone(html string) bool {
	if strings.Contains(html, "See All Buying Options") {
		return true
	}
	return false
}
func lightningDeal(html string) bool {
	if strings.Contains(html, "Lightning Deal") {
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
	doc.Find(".zg_hrsr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find(".zg_hrsr_rank").Text()
		title := s.Find(".zg_hrsr_rank").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
	return result
}

func bestSellersRank(html string) ([]string, []string) {
	var result []string                                                                    //leaderboard
	var rankseller []string                                                                //rank
	rank := regexp.MustCompile("#\\s*([0-9,]+)")                                           //find rank
	rankin1 := regexp.MustCompile("in\\s*([0-9A-Za-z &-]+)\\s*<<a href='/gp/bestsellers/") //find rank in
	rankin2 := regexp.MustCompile(">\\s*([0-9A-Za-z &'-]+)</a></span>")                    //find rank in
	for _, value := range strings.Split(html, "\n") {
		if strings.Contains(value, "<a href='/gp/bestsellers/") { //style 1
			value = strings.TrimSpace(strings.ReplaceAll(value, "(", "<"))
			if len(rankin1.FindStringSubmatch(value)) > 1 {
				if len(rank.FindStringSubmatch(value)) > 1 {
					result = append(result, rankin1.FindStringSubmatch(value)[1])
					rankseller = append(rankseller, rank.FindStringSubmatch(value)[1])
				}
			} else if len(rankin2.FindStringSubmatch(value)) > 1 {
				if len(rank.FindStringSubmatch(value)) > 1 {
					result = append(result, rankin2.FindStringSubmatch(value)[1])
					rankseller = append(rankseller, rank.FindStringSubmatch(value)[1])
				}
			}
		} else if strings.Contains(value, "<a href=\"/gp/bestsellers/") { //style 2
			rankin1style2 := regexp.MustCompile("in\\s*([0-9A-Za-z &-]+)\\s*<<a href=\"/gp/bestsellers/") //find rank in
			rankstyle2 := regexp.MustCompile("<span class=\"zg_hrsr_rank\">#\\s*([0-9,]+)")
			if strings.Contains(value, "#") {
				value = strings.TrimSpace(strings.ReplaceAll(value, "(", "<"))
				if len(rank.FindStringSubmatch(value)) > 1 && len(rankin1style2.FindStringSubmatch(value)) > 1 {
					result = append(result, rankin1style2.FindStringSubmatch(value)[1])
					rankseller = append(rankseller, rank.FindStringSubmatch(value)[1])
				}
			} else if strings.Contains(value, "zg_hrsr_ladder") {
				if len(rankin2.FindStringSubmatch(value)) > 1 && len(rankstyle2.FindStringSubmatch(html)) > 1 {
					result = append(result, rankin2.FindStringSubmatch(value)[1])
					rankseller = append(rankseller, rankstyle2.FindStringSubmatch(html)[1])
				}
			}
		}
	}
	//fmt.Println(result)
	//fmt.Println(rankseller)
	return result, rankseller
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
		cell = row.AddCell()
		cell.Value = "Price"
		cell = row.AddCell()
		cell.Value = value.price
		if value.buyboxgone {
			cell = row.AddCell()
			cell.Value = "buybox gone"
		}
		if value.currentlyunavailable {
			cell = row.AddCell()
			cell.Value = "Currently Unavailable"
		}
		if value.lightningdeal {
			cell = row.AddCell()
			cell.Value = "lightning deal"
		}
		for _, j := range value.otherseller {
			cell = row.AddCell()
			cell.Value = "other seller"
			cell = row.AddCell()
			cell.Value = strings.Join(j, " for ")
		}
		for i, j := range value.rank {
			row = first.AddRow()
			row.SetHeightCM(1)
			cell = row.AddCell()
			cell.Value = value.name
			cell = row.AddCell()
			cell.Value = value.leaderboard[i]
			cell = row.AddCell()
			cell.Value = j
		}
		/*for _, j := range value.otherseller {
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
		}*/
		row = first.AddRow()
		cell = row.AddCell()
		cell.Value = value.name
		cell = row.AddCell()
		cell.Value = "Sold"
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

	content, err := ioutil.ReadFile("url.txt")
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
		if len(findSubURL.FindStringSubmatch(html)) > 1 {
			temp.otherseller, temp.price = otherSeller(findSubURL.FindStringSubmatch(html)[1])
		}
		temp.buyboxgone = buyBoxGone(html)
		temp.lightningdeal = lightningDeal(html)
		temp.leaderboard, temp.rank = bestSellersRank(html)
		temp.currentlyunavailable = currentlyUnavailable(html)
		product = append(product, temp)
		fmt.Println(temp)
		time.Sleep(1000 * time.Millisecond)
	}
	//fmt.Println(product)
	output(product)
}

func main() {
	webCrawler()
	fmt.Println("Finish！")
	fmt.Scanln()
}
