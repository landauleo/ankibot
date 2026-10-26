package main

import (
	"ankibot/internal/client/anki"
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //корень дерева контекстов
	defer cancel()                                                          //гарантия очистки ресурсов контекста

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("❌ Ошибка: переменная окружения TELEGRAM_BOT_TOKEN не задана")
	}

	ankiClient := anki.NewClient("http://127.0.0.1:8765", 2*time.Second) //создали объект-сервис и передали в него строки конфигурации

	if err := ankiClient.EnsureAppRunning(ctx); err != nil {
		log.Fatalf("❌ Не удалось запустить Anki: %v", err)
	}

	decks, err := ankiClient.GetDeckNames(ctx)
	if err != nil {
		log.Fatalf("❌ Ошибка: %v", err)
	}
	fmt.Printf("✅ Успешное подключение! Ваши колоды в Anki: %v\n", decks)

	cardIDs, err := ankiClient.FindCards(ctx, decks[3])
	if err != nil {
		log.Fatalf("❌ Ошибка: %v", err)
	}

	fmt.Printf("✅ Найдено карт в колоде: %d. ID карт: %v\n", len(cardIDs), cardIDs)

	firstCardID := cardIDs[0]
	cardsInfo, err := ankiClient.GetCardsInfo(ctx, []int64{firstCardID})
	if err != nil {
		log.Fatalf("❌ Ошибка при получении данных карты: %v", err)
	}

	card := cardsInfo[0]
	fmt.Printf("🎴 Вопрос (Front): %s\n", card.Fields["Front"].Value)
	fmt.Printf("💡 Ответ (Back): %s\n", card.Fields["Back"].Value)
}
