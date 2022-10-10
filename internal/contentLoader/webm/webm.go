package webm

import (
	"github.com/kapitan123/telegrofler/internal/contentLoader"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
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
	f, err := os.CreateTemp("", "transcode_*.mp4")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()
	cmd := exec.Command("ffmpeg",
		"-i", url, "-y", "-c:v", "libx264", "-preset", "veryfast", f.Name(),
	)
	_, err = cmd.Output()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			log.Warn("ffmpeg exited with error: ", string(exitErr.Stderr))
		}
		return nil, err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return &contentLoader.VideoMeta{
		Id:   url,
		Data: data,
		Type: "mp4",
	}, nil

}
func (e *MetaExtractor) IsServingUrl(url string) bool {
	return strings.HasSuffix(url, ".webm")
}
