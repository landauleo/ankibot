package main

import (
	"ankibot/internal/client/anki"
	"ankibot/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// интерфейсы никогда не передаются по указателю, под капотом передается конкретная реализация,
func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Print(w, "Hello, world")
}
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //корень дерева контекстов
	defer cancel()                                                          //гарантия очистки ресурсов контекста

	ankiClient := anki.NewClient("http://localhost:8999/v1/cards/random", 5*time.Second)
	cardService := service.NewCardService(ankiClient)

	deckName := "espagñol"
	card, err := cardService.GetNextCard(ctx, deckName)
	if err != nil {
		log.Fatalf("Ошибка при получении карточки: %v", err)
	}

	fmt.Printf("Успешно получена карточка: %s\n", card)
}
