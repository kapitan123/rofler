package lovetik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/kapitan123/telegrofler/pkg/source"
)

const (
	noWatermarkType   = "<b> MP4</b> (NO watermark)"
	contentType       = "application/x-www-form-urlencoded; charset=UTF-8"
	sourceLink        = "https://lovetik.com/api/ajax/search"
	sourceType        = "tiktok"
	mobilePrefixRegex = `https:\/\/[a-zA-Z]{2}\.tiktok\.com\/`
)

type LoveTikClient struct {
}

func New() *LoveTikClient {
	ttc := &LoveTikClient{}
	return ttc
}

func (ttc *LoveTikClient) ExtractVideoFromUrl(tikUrl string) (*source.ExtrctedVideoItem, error) {
	escapedUrl := url.QueryEscape(tikUrl)
	jsonStr := []byte("query=" + escapedUrl)
	resp, err := http.Post(sourceLink, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the request has failde with statuscode %d. Data: %s", resp.StatusCode, body)
	}

	sr := &SearchResult{}

	json.Unmarshal(body, sr)

	for _, l := range sr.Links {
		if l.Type != noWatermarkType {
			continue
		}

		log.Print("Found no watermark video: ", l.DownloadAddr)
		b, err := source.DownloadBytesFromUrl(l.DownloadAddr)

		if err != nil {
			return nil, err
		}

		lti := &source.ExtrctedVideoItem{
			Id:      sr.Vid,
			Payload: b,
			Title:   sr.Desc,
			Type:    sourceType,
		}
		return lti, nil
	}

	return nil, nil
}

func IsMatchingUrl(text string) bool {
	r := regexp.MustCompile(mobilePrefixRegex)
	return r.MatchString(text)
}
