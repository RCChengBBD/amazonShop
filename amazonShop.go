package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/NUKBBD/amazonShop/plugins/touchpage"
)

const (
	review_Tracking = "review tracking.txt"
	inputfile       = "url.txt"
	price           = false //For Yawen to grep price
	commentAndStar  = true  //For Lucy To grep Star and comment
	touch           = false //For JJ to touch competitor's page
)

//ReadForKeyWord search for amazon keyword return string array for keyWord
func ReadForKeyWord(file string) []string {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Read file error: ", err)
	}
	body := strings.TrimSpace(string(content))

	var keyWord = []string{}

	for _, key := range strings.Split(body, "\n") {
		keyWord = append(keyWord, strings.TrimSpace(key))
	}
	return keyWord
}

func main() {
	//webCrawler()
	//review()
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()
		if price {
			webCrawler()
		}
	}()

	go func() {
		defer wg.Done()
		if commentAndStar {
			review()
		}
	}()

	go func() {
		defer wg.Done()
		if touch {
			keyword := ReadForKeyWord(inputfile)
			touchpage.TouchKeywordPage(keyword)
		}
	}()

	wg.Wait()
	fmt.Println("Finish！")
	fmt.Scanln()
}
