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
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	req.Header.Set("cookie", "cookie: session-id=131-2340817-4143446; session-id-time=2082787201l; i18n-prefs=USD; ubid-main=130-9943017-9783306; x-wl-uid=1r1GuTQO+ZdaxMTNCQaBPkPjV1JoH/7k62hv/n+PgwdbaOywIv/oT43QJi0BLCdSKhI+FW+34KLA=; sp-cdn=\"L5Z9:TW\"; session-token=I4nDGHJRCv8peqQiV8somyA3CxVNvq8YC58ENj9DnVoNXEHUf4z2eQnJ9OQmXHzZVvVbjgXanYyfYJdeUplNMieHfD6yzkiZcoOwvKr+03vwrFj9i3D96uEunM6XVYYB4Rxz9HcvP8+nDhmYfpxM0kPHYzV5Pe0bKMzorC+AzoGsF8XfBUe8g4cwC/LixoFc; lc-main=zh_TW; csm-hit=tb:s-3SVQ6MVG3K376N91P8YV|1585799864163&t:1585799864378&adb:adblk_yes")
	req.Header.Set("Referer", url)
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
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
	ioutil.WriteFile("otherSeller", []byte(response), 0777)
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
	spiltHtml := strings.Split(html, "\n")
	for n, value := range spiltHtml {
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
				if len(rankin2.FindStringSubmatch(value)) > 1 && len(rankstyle2.FindStringSubmatch(spiltHtml[n-1])) > 1 {
					result = append(result, rankin2.FindStringSubmatch(value)[1])
					rankseller = append(rankseller, rankstyle2.FindStringSubmatch(spiltHtml[n-1])[1])
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
		err = ioutil.WriteFile(temp.name, []byte(html), 0644)
		if err != nil {
			panic(err)
		}
		if len(findSubURL.FindStringSubmatch(html)) > 1 {
			fmt.Println(findSubURL.FindStringSubmatch(html))
			temp.otherseller, temp.price = otherSeller(findSubURL.FindStringSubmatch(html)[1])
		}
		temp.buyboxgone = buyBoxGone(html)
		temp.lightningdeal = lightningDeal(html)
		temp.leaderboard, temp.rank = bestSellersRank(html)
		temp.currentlyunavailable = currentlyUnavailable(html)
		product = append(product, temp)
		fmt.Println(temp)
		time.Sleep(5000 * time.Millisecond)
	}
	//fmt.Println(product)
	output(product)
}

type star struct {
	name             string
	url              string
	onestarurl       string
	twostarurl       string
	threestarurl     string
	productNumber    string
	onestarmessage   []messages
	twostarmessage   []messages
	threestarmessage []messages
}

type messages struct {
	message string
	author  string
	title   string
	date    string
}

func (s *star) getmessage(input string, number string) []messages {
	var message []messages
	context := strings.Split(input, "\n")
	for index, value := range context {
		if strings.Contains(value, "customer_review-") {
			var mess messages
			author := regexp.MustCompile("<span class=\"a-profile-name\">\\s*([^<]+)")
			//m, _ := regexp.Compile("<span>([.]+)</span>")
			date := regexp.MustCompile("Reviewed in the United States on\\s*([0-9,a-zA-Z ]+)")
			if len(author.FindStringSubmatch(value)) > 1 {
				mess.author = author.FindStringSubmatch(value)[1]
				mess.message = context[index+22]
				mess.message = strings.Replace(mess.message, "<span>", "", 1)
				mess.message = strings.Replace(mess.message, "</span>", "", 1)
				mess.message = strings.ReplaceAll(mess.message, "<br>", "\n")
				mess.message = strings.ReplaceAll(mess.message, "<br />", "\n")
				//fmt.Println(mess.author)
			} else {
				fmt.Println(s.name + " " + number + " star can not find author")
				continue
			}
			/*if len(m.FindStringSubmatch(context[index+10])) > 1 {
				mess.title = m.FindStringSubmatch(context[index+10])[1]
			} else {
				fmt.Println(s.name + " " + number + " star can not find title")
			}*/
			if len(date.FindStringSubmatch(context[index+12])) > 1 {
				mess.date = date.FindStringSubmatch(context[index+12])[1]
			} else {
				fmt.Println(s.name + " " + number + " star can not find date")
			}

			//fmt.Println(mess.message)
			message = append(message, mess)
		}
	}
	return message
}

func review() {
	content, err := ioutil.ReadFile("url")
	if err != nil {
		fmt.Println("Read file error: ", err)
	}
	body := string(content)

	parseurl := strings.Split(body, "\n")
	//var threestar []star
	//var twostar []star
	//var onestar []star
	result := []star{}
	productNum, _ := regexp.Compile("dp/\\s*([0-9a-zA-z?=;_&]+)")
	for _, value := range parseurl {
		var temp star
		temp.productNumber = productNum.FindStringSubmatch(value)[1]
		fmt.Println(temp.productNumber)
		temp.name = strings.TrimSpace(strings.Split(value, "\\")[0])
		temp.url = strings.TrimSpace(strings.Split(value, "\\")[1])
		temp.onestarurl = "https://www.amazon.com/product-reviews/" + temp.productNumber + "/ref=acr_dp_hist_1?ie=UTF8&filterByStar=one_star&reviewerType=all_reviews#reviews-filter-bar"
		temp.twostarurl = "https://www.amazon.com/product-reviews/" + temp.productNumber + "/ref=acr_dp_hist_1?ie=UTF8&filterByStar=two_star&reviewerType=all_reviews#reviews-filter-bar"
		temp.threestarurl = "https://www.amazon.com/product-reviews/" + temp.productNumber + "/ref=acr_dp_hist_1?ie=UTF8&filterByStar=three_star&reviewerType=all_reviews#reviews-filter-bar"
		_, context := queryhtmlToString(temp.onestarurl)
		//fmt.Println(context)
		temp.onestarmessage = temp.getmessage(context, "one")
		time.Sleep(2000 * time.Millisecond)

		_, context = queryhtmlToString(temp.twostarurl)
		temp.twostarmessage = temp.getmessage(context, "two")
		time.Sleep(2000 * time.Millisecond)

		_, context = queryhtmlToString(temp.threestarurl)
		temp.threestarmessage = temp.getmessage(context, "three")
		result = append(result, temp)

		time.Sleep(5000 * time.Millisecond)
	}
	fmt.Println("Collext finish!\nGenerating excel file, please wait.")
	reviewoutput(result)
}

func reviewoutput(output []star) {
	file, err := xlsx.OpenFile("format.xlsx")
	if err != nil {
		panic(err)
	}
	first := file.Sheets[0]
	row := first.AddRow()
	row.SetHeightCM(1)
	for _, value := range output {
		fmt.Println()
		for _, one := range value.onestarmessage {
			cell := row.AddCell()
			cell.Value = value.name
			cell = row.AddCell()
			cell.Value = value.url
			cell = row.AddCell()
			cell.Value = one.author
			cell = row.AddCell()
			cell.Value = one.date
			cell = row.AddCell()
			cell.Value = "1"
			cell = row.AddCell()
			cell.Value = one.message
			row = first.AddRow()
			row.SetHeightCM(1)
		}
		for _, one := range value.twostarmessage {
			cell := row.AddCell()
			cell.Value = value.name
			cell = row.AddCell()
			cell.Value = value.url
			cell = row.AddCell()
			cell.Value = one.author
			cell = row.AddCell()
			cell.Value = one.date
			cell = row.AddCell()
			cell.Value = "2"
			cell = row.AddCell()
			cell.Value = one.message
			row = first.AddRow()
			row.SetHeightCM(1)
		}
		for _, one := range value.threestarmessage {
			cell := row.AddCell()
			cell.Value = value.name
			cell = row.AddCell()
			cell.Value = value.url
			cell = row.AddCell()
			cell.Value = one.author
			cell = row.AddCell()
			cell.Value = one.date
			cell = row.AddCell()
			cell.Value = "3"
			cell = row.AddCell()
			cell.Value = one.message
			row = first.AddRow()
			row.SetHeightCM(1)
		}

	}

	err = file.Save("ReviewOutput.xlsx")
	if err != nil {
		panic(err)
	}
}

func main() {
	//webCrawler()
	review()
	fmt.Println("Finish！")
	fmt.Scanln()
}
