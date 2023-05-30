package media

import (
	"regexp"
	"time"

	"github.com/kapitan123/telegrofler/service/bot/domain"
)

type Media struct {
	Url      string
	PostedOn time.Time
	ChatId   int64
	Poster   domain.UserRef
}

var mobileTiktokRegex = regexp.MustCompile(`https:\/\/[a-zA-Z]{2}\.tiktok\.com\/`)
var youtubeShortsRegex = regexp.MustCompile(`https:\/\/youtube.com\/shorts\/.*`)
var instagramReelRegex = regexp.MustCompile(`https:\/\/www.instagram.com\/reel\/.*`)

var supportedMasks = []*regexp.Regexp{
	mobileTiktokRegex,
	youtubeShortsRegex,
	instagramReelRegex,
}

func IsConvertable(url string) bool {
	for _, regex := range supportedMasks {
		if regex.MatchString(url) {
			return true
		}
	}

	return false
}
