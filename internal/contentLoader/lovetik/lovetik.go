package lovetik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/kapitan123/telegrofler/internal/contentLoader"
)

const (
	MP4Type         = "<b> MP4</b>"
	noWatermarkType = "NO watermark"
	contentType     = "application/x-www-form-urlencoded; charset=UTF-8"
	sourceLink      = "https://lovetik.com/api/ajax/search"
	sourceType      = "tiktok"
)

var mobilePrefixRegex = regexp.MustCompile(`https:\/\/[a-zA-Z]{2}\.tiktok\.com\/`)

type Lovetik struct {
}

func New() *Lovetik {
	return &Lovetik{}
}

func (s *Lovetik) ExtractVideoMeta(tikUrl string) (*contentLoader.VideoMeta, error) {
	escapedUrl := url.QueryEscape(tikUrl)
	jsonStr := []byte("query=" + escapedUrl)
	resp, err := http.Post(sourceLink, contentType, bytes.NewBuffer(jsonStr))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the request has failde with statuscode %d. Data: %s", resp.StatusCode, body)
	}

	sr := &SearchResult{}

	json.Unmarshal(body, sr)

	for _, l := range sr.Links {
		if l.WatermarkType != noWatermarkType || l.Type != MP4Type {
			continue
		}

		log.Infof("found video with no watermark: ", l.DownloadAddr)

		meta := &contentLoader.VideoMeta{
			Id:          sr.Vid,
			DownloadUrl: l.DownloadAddr,
			Title:       sr.Desc,
			Type:        sourceType,
		}

		return meta, nil
	}

	return nil, fmt.Errorf("metadata could not be parsed. Probably the format of it was changed")
}

func (s *Lovetik) IsServingUrl(text string) bool {
	return mobilePrefixRegex.MatchString(text)
}
