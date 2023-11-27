package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"TelegramBot/lib/e"
)

// Представляет интерфейс хранилища
type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

// Возвращается, когда нет сохраненных страниц в хранилище
var ErrNoSavedPages = errors.New("нет сохраненных страниц")

// Представляет собой структуру данных страницы с URL и именем пользователя
type Page struct {
	URL      string
	UserName string
}

// Вычисляет хэш-сумму для страницы, используя SHA-1 для URL и имени пользователя
func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("невозможно вычислить хэш", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("невозможно вычислить хэш", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
