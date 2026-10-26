package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	client AnkiClientForBot
}

func NewBot(api *tgbotapi.BotAPI, client AnkiClientForBot) *Bot {
	return &Bot{api: api, client: client}
}

type AnkiClientForBot interface {
	GetRandomCard(ctx context.Context, deckName string) (string, error)
}

func (b *Bot) Start(ctx context.Context) error {
	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60

	updatesChan := b.api.GetUpdatesChan(update)

	for update := range updatesChan {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-updatesChan:
			if update.Message == nil {
				continue
			}
			//TODO
		}
	}
	return nil

}
