package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/NUKBBD/amazonShop/plugins/touchpage"
)

const (
	inputfile      = "url.txt"
	price          = false //For Yawen to grep price
	commentAndStar = false //For Lucy To grep Star and comment
	touch          = true  //For JJ to touch competitor's page
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

	if price {
		webCrawler()
	}
	if commentAndStar {
		review()
	}
	if touch {
		keyword := ReadForKeyWord(inputfile)
		touchpage.TouchKeywordPage(keyword)
	}

	fmt.Println("Finish！")
	fmt.Scanln()
}
