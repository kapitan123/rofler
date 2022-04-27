package source

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

	client := &http.Client{
		Transport: tr,
		Timeout:   50 * time.Second,
	}

	resp, err := client.Get(dUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the request has failde with statuscode %d. Data: %s", resp.StatusCode, body)
	}

	return body, nil
}
