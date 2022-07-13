package lovetik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/kapitan123/telegrofler/internal/contentLoader"
)

const (
	noWatermarkType   = "<b> MP4</b> (NO watermark)"
	contentType       = "application/x-www-form-urlencoded; charset=UTF-8"
	sourceLink        = "https://lovetik.com/api/ajax/search"
	sourceType        = "tiktok"
	mobilePrefixRegex = `https:\/\/[a-zA-Z]{2}\.tiktok\.com\/`
)

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

	io.ReadAll(resp.Body)

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

		meta := &contentLoader.VideoMeta{
			Id:          sr.Vid,
			DownloadUrl: l.DownloadAddr,
			Title:       sr.Desc,
			Type:        sourceType,
		}

		return meta, nil
	}

	return nil, fmt.Errorf("Metadata could not be parsed. Probably the format of it was changed.")
}

func (s *Lovetik) IsServingUrl(text string) bool {
	r := regexp.MustCompile(mobilePrefixRegex)
	return r.MatchString(text)
}
