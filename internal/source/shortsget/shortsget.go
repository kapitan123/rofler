package shortsget

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kapitan123/telegrofler/internal/source"
)

const (
	sourceLink   = "https://shortsget.com/"
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
	// AK TODO make requests in parallel
	vInfo, err := getVideoInfo(youtubeUrl)
	if err != nil {
		return nil, err
	}

	queryParam := "download?videoURL=" + youtubeUrl + "&itag=18&format=mp4"

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

func getVideoInfo(escapedUrl string) (*VideoDetails, error) {
	infoUrl := "videoInfo?videoURL=" + escapedUrl
	resp, err := http.DefaultClient.Get(sourceLink + infoUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the request has failde with statuscode %d. Data: %s", resp.StatusCode, body)
	}

	vidResp := &Response{}
	json.Unmarshal(body, vidResp)

	return &vidResp.VideoDetails, nil
}

func IsMatchingUrl(text string) bool {
	return strings.HasPrefix(text, sourcePrefix)
}

type Response struct {
	VideoDetails VideoDetails `json:"videoDetails"`
	Formats      []Format     `json:"formats"`
}
type VideoDetails struct {
	Title   string `json:"title"`
	VideoId string `json:"videoId"`
}

type Format struct {
	Itag         int    `json:"itag"`
	MimeType     string `json:"mimeType"`     //"mimeType": "video/mp4; codecs=\"avc1.42001E, mp4a.40.2\""
	QualityLabel string `json:"qualityLabel"` //"qualityLabel": "240p"
	Codec        string `json:"codec"`
}
