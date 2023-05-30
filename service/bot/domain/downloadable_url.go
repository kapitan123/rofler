package domain

import "regexp"

type DownloadableUrl struct {
	url string
}

var mobileTiktokRegex = regexp.MustCompile(`https:\/\/[a-zA-Z]{2}\.tiktok\.com\/`)
var youtubeShortsRegex = regexp.MustCompile(`https:\/\/youtube.com\/shorts\/.*`)
var instagramReelRegex = regexp.MustCompile(`https:\/\/www.instagram.com\/reel\/.*`)

var supportedMasks = []*regexp.Regexp{
	mobileTiktokRegex,
	youtubeShortsRegex,
	instagramReelRegex,
}

func TryParseFromSupportedSources(url string) (DownloadableUrl, bool) {
	for _, regex := range supportedMasks {
		if regex.MatchString(url) {
			return DownloadableUrl{url}, true
		}
	}

	return DownloadableUrl{}, false
}
