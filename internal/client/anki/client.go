package anki

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// не нужно писать implements AnkiClient, так как тут implicit interfaces ->
// если у *anki.Client есть метод GetRandomCard, значит, он автоматически реализует интерфейс AnkiClient
type Client struct {
	baseUrl    string
	httpClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseUrl:    baseURL,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) GetRandomCard(ctx context.Context, deckName string) (string, error) {
	//AnkiConnect устроена по спецификации JSON-RPC 2.0 -> все запросы отправляются методом POST
	payload := map[string]interface{}{
		"action":  "findCards",
		"version": 6,
		"params":  map[string]string{"query": "deck:" + deckName},
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseUrl, bytes.NewBuffer(jsonBody))

	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
}
