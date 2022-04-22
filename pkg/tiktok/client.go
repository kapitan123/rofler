// package tiktok

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"regexp"
// 	"strings"

// 	log "github.com/sirupsen/logrus"
// )

// // Reference: https://github.com/drawrowfly/tiktok-scraper
// var defaultHeaders = map[string]string{
// 	"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
// 	"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36",
// 	"referer":    "https://www.tiktok.com/",
// }

// const tikTokRootUrl = "https://www.tiktok.com/"

// type TikTokClient struct {
// 	cookies []http.Cookie
// }

// func New() *TikTokClient {
// 	log.Info("Tiktok client. Loading cookies.")

// 	cookies, err := fetchCookiesFromTikTok()

// 	if err != nil {
// 		log.Error("Cookis fetch has failed future requests might fail: ", err)
// 	} else {
// 		log.Infof("Tiktok client. Cookies loaded.", cookies)
// 	}

// 	ttc := &TikTokClient{cookies}
// 	return ttc
// }

// func (ttc *TikTokClient) DownloadVideoFromItem(item *Item) ([]byte, error) {
// 	return ttc.getFromTikTok(item.Video.DownloadAddr)
// }

// func fetchCookiesFromTikTok() ([]http.Cookie, error) {
// 	log.Info("Loading cookies.")
// 	req, err := http.NewRequest(http.MethodGet, "https://www.tiktok.com/", nil)

// 	if err != nil {
// 		return nil, err
// 	}

// 	for name, h := range defaultHeaders {
// 		req.Header.Set(name, h)
// 	}

// 	resp, err := http.DefaultClient.Do(req)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer resp.Body.Close()

// 	if resp.StatusCode != 200 {
// 		details, _ := ioutil.ReadAll(resp.Body)
// 		return nil, fmt.Errorf("the request has failde with statuscode %d. Data: %s", resp.StatusCode, details)
// 	}

// 	cookies := make([]http.Cookie, 10)
// 	for _, cookie := range resp.Cookies() {
// 		if cookie.Name == "tt_webid_v2" {
// 			continue
// 		}
// 		cookies = append(cookies, *cookie)
// 	}

// 	return cookies, nil
// }

// // Downloads the video item given the username. (@ not included) and video id in the item URL.
// func (ttc *TikTokClient) GetItemByUrl(url string) (*Item, error) {
// 	log.Info("Getting a video item from tikot.")
// 	req, err := http.NewRequest(http.MethodGet, url, nil)

// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, c := range ttc.cookies {
// 		req.AddCookie(&c)
// 	}

// 	cookiess := http.Cookie{
// 		Name:   "sid_tt",
// 		Value:  "1d47ff3138d5f6eac04c46663fa7db76",
// 		Path:   "/",
// 		Domain: ".tiktok.com",
// 	}

// 	req.AddCookie(&cookiess)

// 	for name, h := range defaultHeaders {
// 		req.Header.Set(name, h)
// 	}

// 	log.Error("Built request headers ", req.Header)
// 	log.Error("Built request cookie ", req.Cookies())

// 	// AK TODO change to newClient using cookiejar
// 	client := http.Client{}
// 	resp, err := client.Do(req)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer resp.Body.Close()

// 	if resp.StatusCode != 200 {
// 		details, _ := ioutil.ReadAll(resp.Body)
// 		return nil, fmt.Errorf("the request has failde with statuscode %d. Url: %s. Data: %s", resp.StatusCode, url, details)
// 	}

// 	b, err := ioutil.ReadAll(resp.Body)

// 	if err != nil {
// 		return nil, err
// 	}

// 	stringified := string(b[:])

// 	return ttc.getPersistedDataScript(stringified)
// }

// func (ttc *TikTokClient) getFromTikTok(url string) (b []byte, err error) {
// 	log.Info("Pulling video data from tikot.")
// 	resp, err := ttc.doAuthenticatedRequest(http.MethodGet, url)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer resp.Body.Close()

// 	return ioutil.ReadAll(resp.Body)
// }

// // AK TODO Change to newAuthenticatedRequest
// // Submits a request to tiktok. If used the body should be closed manually.
// func (ttc *TikTokClient) doAuthenticatedRequest(method string, url string) (*http.Response, error) {
// 	req, err := http.NewRequest(method, url, nil)

// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, c := range ttc.cookies {
// 		req.AddCookie(&c)
// 	}

// 	for name, h := range defaultHeaders {
// 		req.Header.Set(name, h)
// 	}

// 	cookiess := http.Cookie{
// 		Name:   "sid_tt",
// 		Value:  "1d47ff3138d5f6eac04c46663fa7db76",
// 		Path:   "/",
// 		Domain: ".tiktok.com",
// 	}
// 	req.AddCookie(&cookiess)

// 	log.Error("Built request headers ", req.Header)
// 	log.Error("Built request cookie ", req.Cookies())

// 	// AK TODO change to newClient using cookiejar

// 	resp, err := http.DefaultClient.Do(req)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if resp.StatusCode != 200 {
// 		log.Error("resp cookies: ", resp.Cookies())
// 		log.Error("resp header: ", resp.Header)
// 		details, _ := ioutil.ReadAll(resp.Body)
// 		return nil, fmt.Errorf("the request has failde with statuscode %d. Url: %s Data: %s", resp.StatusCode, url, details)
// 	}

// 	return resp, nil
// }

// func (ttc *TikTokClient) getPersistedDataScript(content string) (*Item, error) {
// 	r := regexp.MustCompile(`<script id=\"sigi-persisted-data\">[\s\S]*?<\/script>`)

// 	scritpRawText := r.FindString(content)
// 	log.Info("Raw tiktok page: ", scritpRawText)

// 	// The object we need starts with ItemModule
// 	// And ends with the next property UserModule
// 	var startmarker, endmarker = "\"ItemModule\":", ",\"UserModule\":"

// 	// The object inside ItemModule starts with vide id "{ "7075876125477489966":
// 	// We assume that the objetid lenght stays constatnt
// 	idOffset, closingBracketOffest := 23, 1

// 	startIndex := strings.Index(content, startmarker) + len(startmarker) + idOffset
// 	endIndex := strings.Index(content, endmarker) - closingBracketOffest

// 	scritpRawText = content[startIndex:endIndex]

// 	data := &Item{}
// 	err := json.Unmarshal([]byte(scritpRawText), data)

// 	if err != nil {
// 		log.Error(`Error unmarshalling json. The object is malformed.
// 		Probably the page structure has changed.
// 		The raw text after trimming is: `, scritpRawText)

// 		return nil, err
// 	}

// 	return data, nil
// }