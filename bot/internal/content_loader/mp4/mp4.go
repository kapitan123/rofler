package mp4

import (
	"net/http"
	"strings"
	"time"

	contentLoader "github.com/kapitan123/telegrofler/bot/internal/content_loader"
)

type MetaExtractor struct {
	client *http.Client
}

var _ contentLoader.VideoMetaExtractor = (*MetaExtractor)(nil)

func New() *MetaExtractor {
	return &MetaExtractor{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (e *MetaExtractor) ExtractVideoMeta(url string) (*contentLoader.VideoMeta, error) {
	return &contentLoader.VideoMeta{
		Id:          url,
		DownloadUrl: url,
		Type:        "mp4",
	}, nil

}
func (e *MetaExtractor) IsServingUrl(url string) bool {
	return strings.HasSuffix(url, ".mp4")
}
