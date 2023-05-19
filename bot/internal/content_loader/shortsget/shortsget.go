package shortsget

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	contentLoader "github.com/kapitan123/telegrofler/bot/internal/content_loader"
)

const (
	sourceLink   = "https://shortsget.com/"
	sourceType   = "youtubeshorts"
	sourcePrefix = "https://youtube.com/shorts"
)

type Shortsget struct {
}

func New() *Shortsget {
	return &Shortsget{}
}

func (s *Shortsget) ExtractVideoMeta(youtubeUrl string) (*contentLoader.VideoMeta, error) {
	vInfo, err := getVideoInfo(youtubeUrl)
	if err != nil {
		return nil, err
	}

	queryParam := "download?videoURL=" + youtubeUrl + "&itag=18&format=mp4"

	lti := &contentLoader.VideoMeta{
		Id:          vInfo.VideoId,
		DownloadUrl: sourceLink + queryParam,
		Title:       vInfo.Title,
		Type:        sourceType,
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

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the request has failde with statuscode %d. Data: %s", resp.StatusCode, body)
	}

	vidResp := &Response{}
	json.Unmarshal(body, vidResp)

	return &vidResp.VideoDetails, nil
}

func (s *Shortsget) IsServingUrl(text string) bool {
	return strings.HasPrefix(text, sourcePrefix)
}
