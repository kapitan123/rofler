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
type RateLimitedCommand struct {
	userProfiles map[int64]profile
	command      command
}

type profile struct {
	count            int
	last             time.Time
	penaltyExpiresOn time.Time
}

// Can live in a separate package as well
func WithRateLimit(cmd command) *RateLimitedCommand {
	return &RateLimitedCommand{
		command:      cmd,
		userProfiles: map[int64]profile{},
	}
}

func (rl *RateLimitedCommand) ShouldRun(m *tgbotapi.Message) bool {
	if rl.command == nil {
		return false
	}

	return rl.command.ShouldRun(m)
}

func (rl *RateLimitedCommand) Handle(ctx context.Context, m *tgbotapi.Message) error {
	userId, isCommand, now := m.From.ID, m.IsCommand(), time.Now()

	if !isCommand {
		return nil
	}

	prf := rl.updateProfile(userId, now)

	if prf.penaltyExpiresOn.After(now) {
		return fmt.Errorf("user is banned from performing requests %d", userId)
	}

	return rl.command.Handle(ctx, m)
}

func (rl *RateLimitedCommand) updateProfile(userId int64, now time.Time) profile {
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
