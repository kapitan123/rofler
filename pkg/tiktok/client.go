package tiktok

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Reference: https://github.com/drawrowfly/tiktok-scraper
var defaultHeaders = map[string]string{
	"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36",
	"referer":    "https://www.tiktok.com/",
}

const tikTokRootUrl = "https://www.tiktok.com/"

type TikTokClient struct {
	cookies map[string]string
}

func New() *TikTokClient {
	log.Info("Tiktok client. Loading cookies.")

	ttc := &TikTokClient{make(map[string]string)}
	err := ttc.initCookies()

	if err != nil {
		log.Error("Cookis fetch has failed future requests might fail: ", err)
	} else {
		log.Info("Tiktok client. Cookies loaded.")
	}

	return ttc
}

func (ttc *TikTokClient) DownloadVideoFromItem(item *Item) ([]byte, error) {
	return ttc.getFromTikTok(item.Video.DownloadAddr)
}

// Downloads the video item given the username. (@ not included) and video id in the item URL.
func (ttc *TikTokClient) GetItemByUrl(url string) (*Item, error) {
	b, err := ttc.getFromTikTok(url)
	if err != nil {
		return nil, err
	}

	stringified := string(b[:])

	return ttc.getPersistedDataScript(stringified)
}

// Poor mans state store

func (ttc *TikTokClient) initCookies() error {
	resp, err := ttc.doAuthenticatedRequest(http.MethodGet, tikTokRootUrl, defaultHeaders, ttc.cookies)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		ttc.cookies[cookie.Name] = cookie.Value
	}

	return nil
}

func (ttc *TikTokClient) getFromTikTok(url string) (b []byte, err error) {
	resp, err := ttc.doAuthenticatedRequest(http.MethodGet, url, defaultHeaders, ttc.cookies)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// Submits a request to tiktok. If used the body should be closed manually.
func (ttc *TikTokClient) doAuthenticatedRequest(method string, url string, headers map[string]string, cookies map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	for name, value := range ttc.cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	for name, value := range defaultHeaders {
		req.Header.Set(name, value)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the request has failde with statuscode %d. Url: %s", resp.StatusCode, tikTokRootUrl)
	}

	return resp, nil
}

func (ttc *TikTokClient) getPersistedDataScript(content string) (*Item, error) {
	r := regexp.MustCompile(`<script id=\"sigi-persisted-data\">[\s\S]*?<\/script>`)

	scritpRawText := r.FindString(content)
	log.Info("Raw tiktok page: ", scritpRawText)

	// The object we need starts with ItemModule
	// And ends with the next property UserModule
	var startmarker, endmarker = "\"ItemModule\":", ",\"UserModule\":"

	// The object inside ItemModule starts with vide id "{ "7075876125477489966":
	// We assume that the objetid lenght stays constatnt
	idOffset, closingBracketOffest := 23, 1

	startIndex := strings.Index(content, startmarker) + len(startmarker) + idOffset
	endIndex := strings.Index(content, endmarker) - closingBracketOffest

	scritpRawText = content[startIndex:endIndex]

	data := &Item{}
	err := json.Unmarshal([]byte(scritpRawText), data)

	if err != nil {
		log.Error(`Error unmarshalling json. The object is malformed. 
		Probably the page structure has changed.
		The raw text after trimming is: `, scritpRawText)

		return nil, err
	}

	return data, nil
}
