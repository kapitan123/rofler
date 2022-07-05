package sourceFactory

import (
	"github.com/kapitan123/telegrofler/internal/services/downloader"
	"github.com/kapitan123/telegrofler/internal/services/downloader/lovetik"
	"github.com/kapitan123/telegrofler/internal/services/downloader/shortsget"
)

type DataExtractionHandler func(string) (*downloader.ExtrctedVideoItem, error)

// AK TODO fix dependancies and move to the source package
func TryGetExtractor(url string) (DataExtractionHandler, bool) {
	if shortsget.IsMatchingUrl(url) {
		return shortsget.ExtractVideoFromUrl, true
	} else if lovetik.IsMatchingUrl(url) {
		return lovetik.ExtractVideoFromUrl, true
	}
	return nil, false
}
