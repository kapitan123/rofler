package middleware

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	limit             = 1
	penaltyMinutes, _ = time.ParseDuration("3m")
)

type middleware interface {
	Handle(context.Context, *tgbotapi.Message) error
}

type RateLimit struct {
	users map[int64]userRequests
}

type userRequests struct {
	requestCount     int
	lastRequest      time.Time
	penaltyExpiresOn time.Time
}

func New() *RateLimit {
	return &RateLimit{}
}

func (h *RateLimit) Handle(ctx context.Context, m *tgbotapi.Message) error {
	userId, isCommand, now := m.From.ID, m.IsCommand(), time.Now()

	if !isCommand {
		return nil
	}

	var requests userRequests

	if requests, ok := h.users[userId]; ok {
		requests.requestCount += 1
		requests.lastRequest = now

		if requests.requestCount > limit {
			requests.penaltyExpiresOn = requests.lastRequest.Add(penaltyMinutes)
		}

	} else {
		requests = userRequests{1, now, now}
		h.users[userId] = requests
	}

	if requests.penaltyExpiresOn.After(now) {
		return fmt.Errorf("user is banned from performing requests")
	}

	return nil
}
