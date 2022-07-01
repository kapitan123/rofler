package sourceFactory

import (
	"github.com/kapitan123/telegrofler/internal/source"
	"github.com/kapitan123/telegrofler/internal/source/lovetik"
	"github.com/kapitan123/telegrofler/internal/source/shortsget"
)

type Extractor func(string) (*source.ExtrctedVideoItem, error)

// AK TODO fix dependancies and move to the source package
func TryGetExtractor(message string) (Extractor, bool) {
	if shortsget.IsMatchingUrl(message) {
		return shortsget.ExtractVideoFromUrl, true
	} else if lovetik.IsMatchingUrl(message) {
		return lovetik.ExtractVideoFromUrl, true
	}
	return nil, false
}
