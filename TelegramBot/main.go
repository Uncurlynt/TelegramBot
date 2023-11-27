package main

import (
	"context"
	"flag"
	"log"

	tgClient "TelegramBot/clients/telegram"
	event_consumer "TelegramBot/consumer/event-consumer"
	"TelegramBot/events/telegram"
	"TelegramBot/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	// Инициализация подключения к базе данных SQLite
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	// Инициализация хранилища базы данных
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	// Инициализация процессора событий Telegram
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	// Инициализация и запуск консьюмера событий
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

// Возвращает токен для доступа к боту из флагов командной строки
func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
