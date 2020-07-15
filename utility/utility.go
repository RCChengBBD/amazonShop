package utility

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func QueryhtmlToResp(url string) *http.Response {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	req.Header.Add("cookie", "cookie: session-id=131-2340817-4143446; session-id-time=2082787201l; i18n-prefs=USD; ubid-main=130-9943017-9783306; x-wl-uid=1r1GuTQO+ZdaxMTNCQaBPkPjV1JoH/7k62hv/n+PgwdbaOywIv/oT43QJi0BLCdSKhI+FW+34KLA=; sp-cdn=\"L5Z9:TW\"; session-token=I4nDGHJRCv8peqQiV8somyA3CxVNvq8YC58ENj9DnVoNXEHUf4z2eQnJ9OQmXHzZVvVbjgXanYyfYJdeUplNMieHfD6yzkiZcoOwvKr+03vwrFj9i3D96uEunM6XVYYB4Rxz9HcvP8+nDhmYfpxM0kPHYzV5Pe0bKMzorC+AzoGsF8XfBUe8g4cwC/LixoFc; lc-main=zh_TW; csm-hit=tb:s-3SVQ6MVG3K376N91P8YV|1585799864163&t:1585799864378&adb:adblk_yes")
	req.Header.Add("Referer", url)
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
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
