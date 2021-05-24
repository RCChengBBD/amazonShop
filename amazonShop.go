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
	commentAndStar  = false //For Lucy To grep Star and comment
	touch           = true  //For JJ to touch competitor's page
)

//ReadForKeyWord search for amazon keyword return string array for keyWord
func ReadForKeyWord(file string) ([]string, []string) {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Read file error: ", err)
	}
	body := strings.TrimSpace(string(content))

	var keyWord = []string{}
	var avoidsWord = []string{}
	var avoids bool = false

	for _, key := range strings.Split(body, "\n") {
		if strings.Contains(key, "avoids:") {
			avoids = true
			continue
		}

		if avoids {
			avoidsWord = append(avoidsWord, strings.TrimSpace(key))
		} else {
			keyWord = append(keyWord, strings.Replace(strings.TrimSpace(key), " ", "+", -1))
		}
	}
	return keyWord, avoidsWord
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
			keyword, avoidsWord := ReadForKeyWord("keyword.txt")
			err := touchpage.TouchKeywordPage(keyword, avoidsWord)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	wg.Wait()
	fmt.Println("Finish！")
	fmt.Scanln()
}
