package showStats

import (
	"bytes"
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
	"github.com/wcharczuk/go-chart"
)

//type metric string

//const (
//	postsMetric    metric = "posts"
//	ractionsMetric        = "reactions"
//)

const commandName = "showStats"

type ShowStats struct {
	messenger messenger
	storage   postStorage
}

type (
	messenger interface {
		SendImg(chatId int64, img []byte, imgName string, caption string) (int, error)
	}

	postStorage interface {
		GetLastWeekPosts(ctx context.Context) ([]storage.Post, error)
	}
)

func New(messenger messenger, storage postStorage) *ShowStats {
	return &ShowStats{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *ShowStats) Handle(ctx context.Context, m *tgbotapi.Message) error {
	posts, err := h.storage.GetLastWeekPosts(ctx) // Should honour the chat id

	if err != nil {
		return err
	}

	authStats, err := groupPostsByUser(posts)

	if err != nil {
		return err
	}

	lines := splitAuthorsToSeries(authStats)
	graph := chart.Chart{Series: lines}

	// check if i can reuse this approach in my rendering
	// I can probably reuse it in wtermarking, hence avoiding byte array copy
	// the same can be done in dowloading video, I should not allocate for the user
	// user should allocate it for himself
	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buffer)

	if err != nil {
		return err
	}
	_, err = h.messenger.SendImg(m.Chat.ID, buffer.Bytes(), "stats.png", "")

	return err
}

// I can just use a map or create a map subclass
// x axis is total count
// y axis is time
// metric is userPosts
// and reactions

func groupPostsByUser(posts []storage.Post) (map[storage.UserRef][]StatPoint, error) {
	authors := make(map[storage.UserRef][]StatPoint)
	for _, p := range posts {
		var points []StatPoint
		if points, ok := authors[p.UserRef]; ok {
			total := float64(len(points))
			points = append(points, StatPoint{Value: total, Day: p.PostedOn})
		} else {
			points = make([]StatPoint, 0)
		}

		authors[p.UserRef] = points
	}

	return authors, nil
}

func splitAuthorsToSeries(authors map[storage.UserRef][]StatPoint) []chart.Series {
	series := make([]chart.Series, len(authors))

	for a, sps := range authors {
		xSet, ySet := splitCoordinates(sps)

		series = append(series, chart.ContinuousSeries{
			Name:    a.DisplayName,
			XValues: xSet,
			YValues: ySet,
		})
	}

	return series
}

func splitCoordinates(points []StatPoint) (xs []float64, ys []float64) {
	for _, sp := range points {
		ys = append(xs, sp.Value)
		xs = append(xs, sp.FloatDate())
	}
	return
}

func (h *ShowStats) ShouldRun(message *tgbotapi.Message) bool {
	// AK TODO check how command is split and hot to pass arguments
	return message.IsCommand() && message.Command() == commandName
}
