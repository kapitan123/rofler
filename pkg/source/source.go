package source

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
)

type Downloader interface {
	ExtractVideoFromUrl(string) (*ExtrctedVideoItem, error)
}

type ExtrctedVideoItem struct {
	Id      string
	Title   string
	Payload []byte
	Type    string
}

func DownloadBytesFromUrl(dUrl string) ([]byte, error) {
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
