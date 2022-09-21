package command

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

// decorator to limit commands abuse
type RateLimit struct {
	userProfiles map[int64]profile
	command      command
}

type profile struct {
	count            int
	last             time.Time
	penaltyExpiresOn time.Time
}

func WithRateLimit(cmd command) *RateLimit {
	return &RateLimit{
		command:      cmd,
		userProfiles: map[int64]profile{},
	}
}

func (rl *RateLimit) ShouldRun(m *tgbotapi.Message) bool {
	return rl.command.ShouldRun(m)
}

func (rl *RateLimit) Handle(ctx context.Context, m *tgbotapi.Message) error {
	userId, isCommand, now := m.From.ID, m.IsCommand(), time.Now()

	if !isCommand {
		return nil
	}

	prf := rl.updateProfile(userId, now)

	if prf.penaltyExpiresOn.After(now) {
		return fmt.Errorf("user is banned from performing requests")
	}

	return rl.command.Handle(ctx, m)
}

func (rl *RateLimit) updateProfile(userId int64, now time.Time) profile {
	prf, ok := rl.userProfiles[userId]

	if !ok {
		prf = profile{1, now, now}
		rl.userProfiles[userId] = prf
		return prf
	}

	prf.count += 1
	prf.last = now

	if prf.count > limit {
		prf.penaltyExpiresOn = prf.last.Add(penaltyMinutes)
	}

	return prf
}
