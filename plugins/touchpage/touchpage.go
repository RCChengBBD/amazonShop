package touchpage

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/NUKBBD/amazonShop/utility"
)

var url = "https://www.amazon.com/s?language=en_US&k="
var baseUrl = "https://www.amazon.com"
var identityWord = "a-size-base-plus a-color-base a-text-normal"
var reg_keyword = regexp.MustCompile("class=\"a-size-base-plus a-color-base a-text-normal\">\\s*([^<]+)")
var reg_href = regexp.MustCompile("\\s*href=\"\\s*([^<]+)\">")
var top_10 = 10     //attack top 10 competitoy
var attackTimes = 5 //ping page times continusly

//TouchKeywordPage return error
func TouchKeywordPage(keyWords, avoidsWords []string) error {
	for _, keyword := range keyWords {
		top_10 = 10
		_, resps := utility.QueryhtmlToString(url + keyword)
		//fmt.Println(resps)
		resp := strings.Split(resps, "\n")
		for _, context := range resp {
			if strings.Contains(context, identityWord) && !utility.IsAvoid(avoidsWords, context) && top_10 > 0 {
				if len(reg_href.FindStringSubmatch(context)) >= 2 {
					fmt.Println("Attack " + reg_keyword.FindStringSubmatch(context)[1] + " for five times")
				} else {
					continue
				}

				if len(reg_href.FindStringSubmatch(context)) >= 2 {
					for i := 0; i < attackTimes; i++ {
						fmt.Println("Ping " + reg_keyword.FindStringSubmatch(context)[1] + " for " + fmt.Sprint(i+1) + " times")
						utility.QueryhtmlToString(baseUrl + reg_href.FindStringSubmatch(context)[1])
						time.Sleep(3 * time.Second)
					}
					top_10 -= 1
				} else {
					fmt.Println("Can not find the url for " + reg_keyword.FindStringSubmatch(context)[1])
				}
			}
		}
	}
	return nil
}
