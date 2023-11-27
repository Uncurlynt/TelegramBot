package event_consumer

import (
	"log"
	"time"

	"TelegramBot/events"
)

// Consumer представляет собой консьюмер событий
type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

// Создает новый экземпляр Consumer с заданным fetcher, processor и batchSize
func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

// Start начинает цикл обработки событий
func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

// Обрабатывает полученные события, передавая каждое событие processor для дальнейшей обработки
func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
