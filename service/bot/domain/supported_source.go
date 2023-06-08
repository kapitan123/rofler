package domain

import (
	"net/url"
	"regexp"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Source struct {
	regexp     *regexp.Regexp
	preprocess func(url *url.URL)
}

var mobileTiktok = Source{regexp: regexp.MustCompile(`https:\/\/[a-zA-Z]{2}\.tiktok\.com\/`)}
var youtubeShorts = Source{regexp: regexp.MustCompile(`https:\/\/youtube\.com\/shorts\/.*`)}
var instagramReel = Source{regexp: regexp.MustCompile(`https:\/\/www.instagram.com\/reel\/.*`), preprocess: stripQueryParams}
var RedditVideo = Source{regexp: regexp.MustCompile(`https:\/\/www.reddit.com\/r\/.*`), preprocess: stripQueryParams}

var supportedSources = []Source{
	mobileTiktok,
	youtubeShorts,
	instagramReel,
	RedditVideo,
}

func stripQueryParams(url *url.URL) {
	url.RawQuery = ""
}

func canBeDownloaded(urlString string) bool {
	_, found := lo.Find(supportedSources, func(ss Source) bool {
		return ss.regexp.MatchString(urlString)
	})

	return found
}

func getSupportedUrl(urlString string) (*url.URL, error) {
	source, found := lo.Find(supportedSources, func(ss Source) bool {
		return ss.regexp.MatchString(urlString)
	})

	if !found {
		return nil, errors.New("url is not supported")
	}

	parsedUrl, err := url.Parse(urlString)

	if err != nil {
		return nil, errors.Wrap(err, "url is malformed")
	}

	if source.preprocess != nil {
		source.preprocess(parsedUrl)
	}

	return parsedUrl, nil
}
