package lovetik

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const noWatermarkType = "<b> MP4</b> (NO watermark)"

type LoveTikClient struct {
}

func New() *LoveTikClient {
	ttc := &LoveTikClient{}
	return ttc
}

func (ttc *LoveTikClient) DownloadVideoFromUrl(tikUrl string) (*LoveTikItem, error) {
	escapedUrl := url.QueryEscape(tikUrl)
	var jsonStr = []byte("query=" + escapedUrl)
	resp, err := http.Post("https://lovetik.com/api/ajax/search", "application/x-www-form-urlencoded; charset=UTF-8", bytes.NewBuffer(jsonStr))

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
		if l.Type == noWatermarkType {
			log.Print("Found watermark video: ", l.DownloadAddr)
			b, err := download(l.DownloadAddr)

			if err != nil {
				return nil, err
			}

			lti := &LoveTikItem{
				Id:      sr.Vid,
				Payload: b,
				Title:   sr.Desc,
			}
			return lti, nil
		}
	}

	return nil, nil
}

func download(dUrl string) ([]byte, error) {

	// AK TODO temp solutions
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // <--- Problem
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(dUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

type SearchResult struct {
	Author string `json:"author"`
	Links  []Link `json:"links"`
	Vid    string `json:"vid"`
	Desc   string `json:"desc"`
}

type Link struct {
	DownloadAddr string `json:"a"`
	Type         string `json:"t"`
}

type LoveTikItem struct {
	Id      string
	Title   string
	Payload []byte
}
