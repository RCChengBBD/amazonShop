package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/NUKBBD/amazonShop/utility"
	"github.com/tealeg/xlsx"
)

type star struct {
	name             string
	url              string
	totalStar        string
	onestarurl       string
	twostarurl       string
	threestarurl     string
	productNumber    string
	onestarmessage   []messages
	twostarmessage   []messages
	threestarmessage []messages
}

type messages struct {
	name      string
	url       string
	message   string
	author    string
	title     string
	date      string
	star      string
	timestamp int64
}

var outputmessage []messages
var emptymessage []string

func (s *star) getmessage(input string, number string) []messages {
	var message []messages
	var url string
	context := strings.Split(input, "\n")
	for index, value := range context {
		if strings.Contains(value, "customer_review-") {
			var mess messages
			reviewIndex := regexp.MustCompile("customer_review-\\s*([A-Za-z0-9]+)")
			author := regexp.MustCompile("<span class=\"a-profile-name\">\\s*([^<]+)")
			date := regexp.MustCompile("Reviewed in the [a-zA-Z ]+ on\\s*([0-9,a-zA-Z ]+)")
			if len(author.FindStringSubmatch(value)) > 1 {
				mess.author = author.FindStringSubmatch(value)[1]
				mess.message = context[index+25]
				mess.message = strings.Replace(mess.message, "<span>", "", 1)
				mess.message = strings.Replace(mess.message, "</span>", "", 1)
				mess.message = strings.ReplaceAll(mess.message, "<br>", "\n")
				mess.message = strings.ReplaceAll(mess.message, "<br />", "\n")
				//fmt.Println(mess.author)
				//fmt.Println(mess.message)
			} else {
				if len(author.FindStringSubmatch(context[index+2])) > 1 {
					mess.author = author.FindStringSubmatch(context[index+2])[1]
					mess.message = context[index+28]
					mess.message = strings.Replace(mess.message, "<span>", "", 1)
					mess.message = strings.Replace(mess.message, "</span>", "", 1)
					mess.message = strings.ReplaceAll(mess.message, "<br>", "\n")
					mess.message = strings.ReplaceAll(mess.message, "<br />", "\n")
					//fmt.Println(mess.author)
					//fmt.Println(mess.message)
				} else {
					fmt.Println(s.name + " " + number + " star can not find author")
					//fmt.Println(value)
					//fmt.Println(index)
					//fmt.Println(s.twostarurl)
					continue
				}
			}
			/*if len(m.FindStringSubmatch(context[index+10])) > 1 {
				mess.title = m.FindStringSubmatch(context[index+10])[1]
			} else {
				fmt.Println(s.name + " " + number + " star can not find title")
			}*/
			if len(date.FindStringSubmatch(context[index+12])) > 1 {
				format := "January 2, 2006"
				mess.date = date.FindStringSubmatch(context[index+12])[1]
				d, err := time.Parse(format, date.FindStringSubmatch(context[index+12])[1])
				mess.timestamp = d.Unix()
				if err != nil {
					fmt.Println(s.name, number, "star date parse error: ", err)
				}
			} else {
				fmt.Println(s.name + " " + number + " star can not find date")
			}

			if len(reviewIndex.FindStringSubmatch(value)) > 1 {
				url = "https://www.amazon.com/gp/customer-reviews/" + reviewIndex.FindStringSubmatch(value)[1] + "/ref=cm_cr_arp_d_rvw_ttl?ie=UTF8"
			} else {
				url = "Can not find url"
			}

			//fmt.Println(mess.message)
			mess.star = number
			mess.name = s.name
			mess.url = url
			message = append(message, mess)
			outputmessage = append(outputmessage, mess)
		}
	}

	if len(message) == 0 {
		fmt.Println(s.name, number, "star comment is empty. Please check.")
		emptymessage = append(emptymessage, s.name)
	}
	return message
}
func reviewoutput(output []messages) {
	file, err := xlsx.OpenFile("format.xlsx")
	if err != nil {
		panic(err)
	}
	first := file.Sheets[0]
	row := first.AddRow()
	row.SetHeightCM(1)
	//fmt.Println(len(output))
	for _, value := range output {
		//fmt.Println()
		cell := row.AddCell()
		cell.Value = value.name
		cell = row.AddCell()
		cell.Value = value.url
		cell = row.AddCell()
		cell.Value = value.author
		cell = row.AddCell()
		cell.Value = value.date
		cell = row.AddCell()
		cell.Value = value.star
		cell = row.AddCell()
		cell.Value = value.message
		row = first.AddRow()
		row.SetHeightCM(1)

	}

	err = file.Save("ReviewOutput.xlsx")
	if err != nil {
		panic(err)
	}
}

func starOutput(output []star) {
	file, err := xlsx.OpenFile("format.xlsx")
	if err != nil {
		panic(err)
	}
	first := file.Sheets[0]
	row := first.AddRow()
	row.SetHeightCM(1)
	for _, value := range output {
		//fmt.Println()
		cell := row.AddCell()
		cell.Value = value.name
		cell = row.AddCell()
		if strings.TrimSpace(value.totalStar) == "" {
			cell.Value = "N/A"
		} else {
			cell.Value = value.totalStar
		}

		row = first.AddRow()
		row.SetHeightCM(1)

	}

	err = file.Save("StarOutput.xlsx")
	if err != nil {
		panic(err)
	}
}
func getPage(input string) int {
	showPage := regexp.MustCompile("class=\"a-size-base\">Showing 1-10 of\\s*([0-9]+)")
	//context := strings.Split(input, "\n")

	if len(showPage.FindStringSubmatch(input)) > 1 {
		messageNumber := showPage.FindStringSubmatch(input)[1]
		result, err := strconv.Atoi(messageNumber)
		if err != nil {
			fmt.Printf("Get message number error %v", err)
		}
		if result%10 == 0 {
			return (result / 10)
		} else {
			return (result / 10) + 1
		}
	}
	return 1
}

type ByTimastamp []messages

func (a ByTimastamp) Len() int           { return len(a) }
func (a ByTimastamp) Less(i, j int) bool { return a[i].timestamp > a[j].timestamp }
func (a ByTimastamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (s *star) setURL() {
	u := strings.Split(s.url, "/ref=")
	s.onestarurl = u[0] + "/ref=cm_cr_unknown?formatType=current_format&language=en_US&reviewerType=all_reviews&filterByStar=one_star&sortBy=recent"
	s.twostarurl = u[0] + "/ref=cm_cr_unknown?formatType=current_format&language=en_US&reviewerType=all_reviews&filterByStar=two_star&sortBy=recent"
	s.threestarurl = u[0] + "/ref=cm_cr_unknown?formatType=current_format&language=en_US&reviewerType=all_reviews&filterByStar=three_star&sortBy=recent"
}

func review() {
	content, err := ioutil.ReadFile(inputfile)
	if err != nil {
		fmt.Println("Read file error: ", err)
	}
	body := strings.TrimSpace(string(content))

	parseurl := strings.Split(body, "\n")
	//var threestar []star
	//var twostar []star
	//var onestar []star
	result := []star{}
	productNum, _ := regexp.Compile("dp/\\s*([0-9a-zA-Z]+)")
	conreview, _ := regexp.Compile("product-reviews/\\s*([0-9a-zA-Z]+)")
	var wg sync.WaitGroup
	//var lock sync.Mutex
	for _, value := range parseurl {
		wg.Add(1)
		go func(value string) {
			defer wg.Done()
			var temp star
			temp.name = strings.TrimSpace(strings.Split(value, "\\")[0])
			temp.url = strings.TrimSpace(strings.Split(value, "\\")[1])
			if len(productNum.FindStringSubmatch(value)) > 1 {
				temp.productNumber = productNum.FindStringSubmatch(value)[1]
			} else {
				if len(conreview.FindStringSubmatch(value)) > 1 {
					temp.productNumber = conreview.FindStringSubmatch(value)[1]
				} else {
					fmt.Println(temp.name + " can not parse url, please confirm.")
				}
			}

			fmt.Println(temp.name)

			temp.setURL()
			_, context := utility.QueryhtmlToString(temp.onestarurl)

			temp.getTotalStar(context) //get the total star

			page := getPage(context)

			urltemp := temp.onestarurl
			for i := 1; i <= page; i++ {

				temp.onestarmessage = temp.getmessage(context, "1")
				urltemp = temp.onestarurl + "&pageNumber=" + strconv.Itoa(i+1)
				_, context = utility.QueryhtmlToString(urltemp)
			}

			urltemp = temp.twostarurl
			_, context = utility.QueryhtmlToString(temp.twostarurl)
			page = getPage(context)
			for i := 1; i <= page; i++ {

				temp.twostarmessage = temp.getmessage(context, "2")
				urltemp := temp.twostarurl + "&pageNumber=" + strconv.Itoa(i+1)
				_, context = utility.QueryhtmlToString(urltemp)
			}

			urltemp = temp.threestarurl
			_, context = utility.QueryhtmlToString(temp.threestarurl)
			page = getPage(context)
			for i := 1; i <= page; i++ {
				temp.threestarmessage = temp.getmessage(context, "3")
				urltemp := temp.threestarurl + "&pageNumber=" + strconv.Itoa(i+1)
				_, context = utility.QueryhtmlToString(urltemp)
			}

			result = append(result, temp)

		}(value)

		//time.Sleep(5000 * time.Millisecond)
	}
	wg.Wait()
	fmt.Println("Collext finish!\nGenerating excel file, please wait.")
	sort.Sort(ByTimastamp(outputmessage))
	reviewoutput(outputmessage)
	starOutput(result)
	fmt.Println("The product of empty message have:")
	/*for _, empty := range emptymessage {
		fmt.Println(empty)
	}*/
}

func (s *star) getTotalStar(input string) {

	context := strings.Split(input, "\n")
	for _, value := range context {
		if strings.Contains(value, "reviewNumericalSummary") || strings.Contains(value, "averageStarRatingNumerical") {
			totalStar := regexp.MustCompile("<span class=\"a-icon-alt\">\\s*([0-9.]+)")
			if len(totalStar.FindStringSubmatch(value)) > 1 {
				if strings.TrimSpace(totalStar.FindStringSubmatch(value)[1]) == "" {
					s.totalStar = "N/A"
				} else {
					s.totalStar = totalStar.FindStringSubmatch(value)[1]
					//fmt.Println(s.totalStar)
				}
				break
			} else {
				s.totalStar = "N/A"
				fmt.Printf("Can not find %s tital star", s.name)
				break
			}
		}
	}
}
