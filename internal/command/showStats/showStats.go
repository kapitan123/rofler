package showStats

import (
	"bytes"
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
	"github.com/wcharczuk/go-chart"
)

type metric string

const (
	postsMetric    metric = "posts"
	ractionsMetric        = "reactions"
)

const commandName = "showStats"

type ShowStats struct {
	messenger messenger
	storage   postStorage
}

type (
	messenger interface {
		SendImg(chatId int64, img []byte, imgName string, caption string) error
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

	stats, err := groupPostsByUser(posts)

	xSet, ySet := splitCoordinates(stats)

	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: xSet,
				YValues: ySet,
			},
		},
	}

	// check if i can reuse this approach in my rendering
	// I can probably reuse it in wtermarking, hence avoiding byte array copy
	// the same can be done in dowloading video, I should not allocate for the user
	// user should allocate it for himself
	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buffer)

	if err != nil {
		return err
	}

	return h.messenger.SendImg(m.Chat.ID, buffer.Bytes(), "stats.png", "")
}

// I can just use a map or create a map subclass
// x axis is total count
// y axis is time
// metric is userPosts
// and reactions

func groupPostsByUser(posts []storage.Post) ([]StatPoint, error) {
	return nil, nil
}

func splitCoordinates(points []StatPoint) ([]float64, []float64) {
	return nil, nil
}

func (h *ShowStats) ShouldRun(message *tgbotapi.Message) bool {
	// AK TODO check how command is split and hot to pass arguments
	return message.IsCommand() && message.Command() == commandName
}
