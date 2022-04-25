package shortsget

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/kapitan123/telegrofler/pkg/source"
)

const (
	sourceLink   = "https://shortsget.com/download"
	sourceType   = "youtubeshorts"
	sourcePrefix = "https://youtube.com/shorts"
)

type ShortsGetClient struct {
}

func New() *ShortsGetClient {
	sgc := &ShortsGetClient{}
	return sgc
}

func (ttc *ShortsGetClient) ExtractVideoFromUrl(youtubeUrl string) (*source.ExtrctedVideoItem, error) {
	escapedUrl := url.QueryEscape(youtubeUrl)

	// AK TODO make requests in parallel
	vInfo, err := getVideoInfo(escapedUrl)
	if err != nil {
		return nil, err
	}

	queryParam := "download?videoURL==" + escapedUrl

	b, err := source.DownloadBytesFromUrl(sourceLink + queryParam)
	if err != nil {
		return nil, err
	}

	lti := &source.ExtrctedVideoItem{
		Id:      vInfo.VideoId,
		Payload: b,
		Title:   vInfo.Title,
		Type:    sourceType,
	}
	return lti, nil
}

func getVideoInfo(escapedUrl string) (*VideoInfo, error) {
	infoUrl := "videoInfo?videoURL==" + escapedUrl
	resp, err := http.DefaultClient.Get(sourceLink + infoUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the request has failde with statuscode %d. Data: %s", resp.StatusCode, body)
	}

	vInfo := &VideoInfo{}
	json.Unmarshal(body, vInfo)
	return vInfo, nil
}

func IsMatchingUrl(text string) bool {
	return strings.HasPrefix(text, sourcePrefix)
}

type VideoInfo struct {
	Title   string `json:"title"`
	VideoId string `json:"videoId"`
}
