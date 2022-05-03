package sourceFactory

import (
	"github.com/kapitan123/telegrofler/internal/source"
	"github.com/kapitan123/telegrofler/internal/source/lovetik"
	"github.com/kapitan123/telegrofler/internal/source/shortsget"
)

func TryGetSource(message string) (source.Downloader, bool) {
	if shortsget.IsMatchingUrl(message) {
		return shortsget.New(), true
	} else if lovetik.IsMatchingUrl(message) {
		return lovetik.New(), true
	}
	return nil, false
}
