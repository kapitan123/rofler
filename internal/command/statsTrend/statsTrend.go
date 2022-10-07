package statsTrend

import (
	"bytes"
	"context"
	"io"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
	"github.com/wcharczuk/go-chart"
)

//type metric string

//const (
//	postsMetric    metric = "posts"
//	ractionsMetric        = "reactions"
//)

const commandName = "statstrend"

type ShowStats struct {
	messenger messenger
	storage   postStorage
}

type (
	messenger interface {
		SendImg(chatId int64, img io.Reader, imgName string, caption string) (int, error)
	}

	postStorage interface {
		GetLastWeekPosts(ctx context.Context, chatid int64) ([]storage.Post, error)
	}
)

func New(messenger messenger, storage postStorage) *ShowStats {
	return &ShowStats{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *ShowStats) Handle(ctx context.Context, m *tgbotapi.Message) error {
	posts, err := h.storage.GetLastWeekPosts(ctx, m.Chat.ID)

	if err != nil {
		return err
	}

	authorStats, err := groupPostsByUser(posts)

	if err != nil {
		return err
	}

	lines := splitAuthorsToSeries(authorStats)

	graph := chart.Chart{
		Title:  "Week Rofler stats",
		Series: lines,
		XAxis: chart.XAxis{
			Name: "postedon",
		},

		YAxis: chart.YAxis{
			Name: "posts",
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buffer)

	if err != nil {
		return err
	}

	_, err = h.messenger.SendImg(m.Chat.ID, buffer, "stats.png", "")

	return err
}

func groupPostsByUser(posts []storage.Post) (map[storage.UserRef][]StatPoint, error) {
	authors := make(map[storage.UserRef][]StatPoint)
	startOfTheWeek := time.Now().AddDate(0, 0, -7)

	for _, p := range posts {
		if points, ok := authors[p.UserRef]; ok {
			total := float64(len(points))
			points = append(points, StatPoint{Value: total, Day: p.PostedOn})
			authors[p.UserRef] = points
		} else {
			authors[p.UserRef] = []StatPoint{{Value: 0, Day: startOfTheWeek}, {Value: 1, Day: p.PostedOn}}
		}
	}

	return authors, nil
}

func splitAuthorsToSeries(authors map[storage.UserRef][]StatPoint) []chart.Series {
	series := make([]chart.Series, len(authors)*2)
	i := 0
	for a, sps := range authors {
		totals, timestamps := splitCoordinates(sps)

		series[i] = chart.TimeSeries{
			Name:    a.DisplayName,
			XValues: timestamps,
			YValues: totals,
		}

		series[i+1] = chart.AnnotationSeries{
			Annotations: []chart.Value2{
				{XValue: 1.0, YValue: 1.0, Label: "One"},
				{XValue: 2.0, YValue: 2.0, Label: "Two"},
				{XValue: 3.0, YValue: 3.0, Label: "Three"},
				{XValue: 4.0, YValue: 4.0, Label: "Four"},
				{XValue: 5.0, YValue: 5.0, Label: "Five"},
			},
		}

		i = i + 2
	}

	return series
}

func splitCoordinates(points []StatPoint) (totals []float64, timestamps []time.Time) {
	for _, sp := range points {
		totals = append(totals, sp.Value)
		timestamps = append(timestamps, sp.Day)
	}
	return
}

func (h *ShowStats) ShouldRun(message *tgbotapi.Message) bool {
	return message.IsCommand() && message.Command() == commandName
}
