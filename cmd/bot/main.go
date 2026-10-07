package main

import (
	"ankibot/internal/client/anki"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //корень дерева контекстов
	defer cancel()                                                          //гарантия очистки ресурсов контекста

	ankiClient := anki.NewClient("http://127.0.0.1:8765", 2*time.Second)

	if err := ankiClient.EnsureAppRunning(ctx); err != nil {
		log.Fatalf("❌ Не удалось запустить Anki: %v", err)
	}

	decks, err := ankiClient.GetDeckNames(ctx)
	if err != nil {
		log.Fatalf("❌ Ошибка: %v", err)
	}
	fmt.Printf("✅ Успешное подключение! Ваши колоды в Anki: %v\n", decks)

	cardIDs, err := ankiClient.FindCards(ctx, decks[0])
	if err != nil {
		log.Fatalf("❌ Ошибка: %v", err)
	}

	fmt.Printf("✅ Найдено карт в колоде: %d. ID карт: %v\n", len(cardIDs), cardIDs)
}
