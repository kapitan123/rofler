package sourceFactory

import (
	"github.com/kapitan123/telegrofler/internal/source"
	"github.com/kapitan123/telegrofler/internal/source/lovetik"
	"github.com/kapitan123/telegrofler/internal/source/shortsget"
)

func TryGetSourceExtractor(message string) (func(string) (*source.ExtrctedVideoItem, error), bool) {
	if shortsget.IsMatchingUrl(message) {
		return shortsget.ExtractVideoFromUrl, true
	} else if lovetik.IsMatchingUrl(message) {
		return lovetik.ExtractVideoFromUrl, true
	}
	return nil, false
}
