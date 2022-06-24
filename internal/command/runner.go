package command

import (
	"context"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type command interface {
	Handle(context.Context, *tgbotapi.Message) error
	ShouldRun(*tgbotapi.Message) bool
}

type Runner struct {
	commands   []command
	numWorkers int
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         *sync.WaitGroup
	ch         chan *tgbotapi.Message
}

func NewRunner(numWorkers int, commands ...command) *Runner {
	return &Runner{
		commands:   commands,
		numWorkers: numWorkers,
		wg:         &sync.WaitGroup{},
		ch:         make(chan *tgbotapi.Message),
	}
}

func (r *Runner) Start(ctx context.Context) {
	r.ctx, r.cancelFunc = context.WithCancel(ctx)
	for i := 0; i < r.numWorkers; i++ {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			for {
				select {
				case <-r.ctx.Done():
					return
				case msg := <-r.ch:
					err := r.run(r.ctx, msg)
					if err != nil {
						log.WithError(err).Error("failed to handle message")
					}
				}
			}
		}()
	}
}

func (r *Runner) Run(ctx context.Context, message *tgbotapi.Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.ctx.Done():
		return r.ctx.Err()
	case r.ch <- message:
		return nil
	}
}

func (r *Runner) run(ctx context.Context, message *tgbotapi.Message) error {
	for _, cmd := range r.commands {
		if cmd.ShouldRun(message) {
			log.Infof("Handled by: %T", cmd)
			return cmd.Handle(ctx, message)
		}
	}
	return nil
}

func (r *Runner) Stop() {
	r.cancelFunc()
	r.wg.Wait()
	close(r.ch)
}
