package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/NUKBBD/amazonShop/utility"
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

func otherSeller(url string) ([][]string, string) { //return seller who is not Amazon,Primo Super-store,KickBOT and SPORTSBOT and the cheapest price.
	var result [][]string

	//var cheapest string
	var cheap float64
	findotherSeller, _ := regexp.Compile("from seller\\s*([a-zA-Z0-9 $.-]+)")
	//findotherPrice, _ := regexp.Compile("and price $\\s*([0-9.]+) <")
	url = "https://www.amazon.com/gp/offer-listing/" + url
	//fmt.Println(url)
	_, response := utility.QueryhtmlToString(url)
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

	content, err := ioutil.ReadFile(inputfile)
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
		_, html := utility.QueryhtmlToString(temp.url)
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

	output(product)
}
