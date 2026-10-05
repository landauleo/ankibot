package service

import (
	"context"
)

// Consumer-defined interface - интерфейс определяется там, где используется
type AnkiClient interface {
	GetRandomCard(ctx context.Context, deckName string) (string, error)
}

// для тестирования и гибкости принято зависеть от интерфейса
type CardService struct {
	ankiClient AnkiClient
}

func NewCardService(ankiClient AnkiClient) *CardService {
	return &CardService{
		ankiClient: ankiClient,
	}
}

func (s *CardService) GetNextCard(ctx context.Context, deckName string) (string, error) {
	card, err := s.ankiClient.GetRandomCard(ctx, deckName)
	if err != nil {
		return "", err
	}

	return card, nil
}
