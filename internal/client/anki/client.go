package anki

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
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

type RequestDto struct {
	Action  string      `json:"action"`
	Version int         `json:"version"`
	Params  interface{} `json:"params,omitempty"` //omitempty - не добавлять в конечный json, если поле nil
}

type ResponseDto struct {
	Result json.RawMessage `json:"result"` //один универсальный ответ для всех запросов AnkiConnect
	Error  *string         `json:"error"`  //ошибка может быть null, но string быть nil не может (максимум "") ->
	//указатель помогает принимать значения nil, если в ответе будет null
}

func (c *Client) GetDeckNames(ctx context.Context) ([]string, error) {
	rawResult, err := c.execute(ctx, "deckNames", nil)
	if err != nil {
		return nil, err
	}

	var deckNames []string
	if err := json.Unmarshal(rawResult, &deckNames); err != nil {
		return nil, fmt.Errorf("ошибка парсинга колод: %w", err)
	}

	return deckNames, nil
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
	return "", nil //TODO
}

func (c *Client) EnsureAppRunning(ctx context.Context) error {
	_, err := c.GetDeckNames(ctx)
	if err == nil {
		return nil
	}

	fmt.Println("🤖 Anki не запущен. Автоматически запускаем /Applications/Anki.app...")

	err = exec.CommandContext(ctx, "open", "-a", "Anki").Run()
	if err != nil {
		return fmt.Errorf("ошибка при попытке запустить Anki.app: %w", err)
	}

	//ждем пока Анки стартанет
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		_, err := c.GetDeckNames(ctx)
		if err == nil {
			fmt.Println("✅ Anki успешно запущен и готов к работе!")
			return nil
		}
	}
	return fmt.Errorf("Anki.app запущен, но AnkiConnect не ответил в течение 10 секунд")
}

func (c *Client) FindCards(ctx context.Context, deckName string) ([]int64, error) {
	params := map[string]string{
		"query": "deck:" + deckName,
	}

	rawResult, err := c.execute(ctx, "findCards", params)
	if err != nil {
		return nil, err
	}
	var cardsIds []int64
	if err := json.Unmarshal(rawResult, &cardsIds); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ID карт: %w", err)
	}

	return cardsIds, nil
}

func (c *Client) execute(ctx context.Context, action string, params interface{}) (json.RawMessage, error) {
	reqBody := RequestDto{
		Action:  action,
		Version: 6,
		Params:  params,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации запроса: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()

	var respDto ResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&respDto); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	if respDto.Error != nil {
		return nil, fmt.Errorf("ошибка AnkiConnect: %s", *respDto.Error)
	}

	return respDto.Result, nil
}
