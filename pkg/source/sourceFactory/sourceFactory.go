package sourceFactory

import (
	"github.com/kapitan123/telegrofler/pkg/source"
	"github.com/kapitan123/telegrofler/pkg/source/lovetik"
	"github.com/kapitan123/telegrofler/pkg/source/shortsget"
)

func TryGetSource(message string) (source.Downloader, bool) {
	if lovetik.IsMatchingUrl(message) {
		return shortsget.New(), true
	} else if shortsget.IsMatchingUrl(message) {
		return lovetik.New(), true
	}
	return nil, false
}
