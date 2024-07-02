package app

import (
	"context"
	"io"
	"net/http"

	"github.com/kapitan123/telegrofler/common/logs"
)

// type Server struct {
// 	fileBucket   *fileBucket
// 	downloader   *downloader
// 	successTopic *successTopic
// }

// AK TODO maybe I don't need interfaces here
type successTopic interface {
	PublishSuccess(ctx context.Context, savedVideoAddr string, originalUrl string) error
}

type fileBucket interface {
	Save(ctx context.Context, fromReader io.Reader) (string, error)
}

type downloader interface {
	DownloadFromUrl(ctx context.Context, url string, out io.Writer) error
}

func NewServer(fileBucket fileBucket, downloader downloader, successTopic successTopic) http.Handler {
	mux := http.NewServeMux()

	addRoutes(
		mux,
		fileBucket,
		downloader,
		successTopic,
	)

	var handler http.Handler = mux

	handler = logs.NewLoggingMiddleware(handler)
	handler = logs.NewTraceIdMiddleware(handler)

	return handler
}

func addRoutes(mux *http.ServeMux, fileBucket fileBucket, downloader downloader, successTopic successTopic) http.Handler {

}
