package utility

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func QueryhtmlToResp(url string) *http.Response {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Host", "www.amazon.com")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Postman-Token", "11b85884-98ae-504a-5ee8-b4b4e1cb2264")
	//req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	//req.Header.Add("cookie", "session-id=130-2503036-6347030; sp-cdn=\"L5Z9:TW\"; ubid-main=132-0247640-4911038; x-main=\"yyP?0IqBWzd1rw2loRPJVjA7ORJvRjZ9WAsUzVqctr5cJCwlLKIo5OcqTE7GUzUN\"; at-main=Atza|IwEBIPrdT388YK_BwvKm_FJ99pgrNnlQn8VSgyjE2tql8AYPioImAyk-G080aYCKoKBKtfOdxnwzm1x_wPYmxkCaxzTd2fYRGtnbIj3QoZO_WC_hHoeiF6zIVea35rDlCK90pQFveMmxvO1WdHW5PLgr_MAbNi75XF_dHIO_O8BAgN5RlXqRdqPYpsXSb-_gLsC8s9FWFIcFh57oFtmak0NtJeGr; sess-at-main=\"L1dyO3msQPfVx3VBRcgY7CcOPNYqA6t+fTeqR4dD7Ko=\"; sst-main=Sst1|PQE1XuvvFybMHRMOjNvpQlbHCZ6T1awFrSCCX05fqBXq6hz-KfLO9JgXRxp8isMKUgP_viQ5iwD0In4uDlLyWKQ5KYDpU9kwpc40-uQvBy13JY_zwRgmj9FocvhTIM69rl166aBkQ05Y_MBT1gBg1cfi4eV8_9wTO8HaSO-dk4FknjMgm-U89l2MEaUHyg_lYxX-w8kLiRlfktSYa1SpX7yLpcVoNExB5lr5EudQH-Zh879ey1N0SogpiLY1aleUnAGdsXzD6XsbXe2kbSkEn4-eTc4P9JNOJjLtVN6KtXVJtGQ; session-id-time=2082787201l; i18n-prefs=TWD; session-token=\"ZePkGX4BjDygiL9WPI37MQ/G0DjcIYYZsg5H2HX60UpGiYQ8l+ik/IH7W2+i1jwSbMTghngSESwmey2+OfLkBqbAKUWQ+vGRCyF2EQgLEEFuIG8Q7iLyqX/0JxlNP7VbA1yXu4C/eBJ4oE+lCoyYArQV7G/k/cjH6gmbcCYHHSmwdnCmU80cjdPNYTAf9NvBaS/04Hcl+xC3/TMvARkFfg==\"; lc-main=en_US; csm-hit=tb:s-8CSRDMYDNWPM0SY3QYZP|1621701261044&t:1621701262120&adb:adblk_yes")
	//req.Header.Add("Referer", url)
	//req.Header.Add("Cache-Control", "no-cache")
	//req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return nil
	}
	/*defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)*/

	return resp
}

func QueryhtmlToString(url string) (*http.Response, string) {

	resp := QueryhtmlToResp(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read all", err)
	}

	return resp, string(body)
}

//return true if context contain the avoids word
//return false if not
func IsAvoid(avoids []string, context string) bool {
	for _, avoid := range avoids {
		if strings.Contains(context, avoid) {
			return true
		}
	}
	return false
}
