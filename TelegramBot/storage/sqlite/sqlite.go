package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"TelegramBot/storage"
)

// Представляет собой хранилище данных, реализованное с использованием SQLite
type Storage struct {
	db *sql.DB
}

// Создает новое SQLite-хранилище
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("невозможно открыть базу данных: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("невозможно подключиться к базе данных: %w", err)
	}

	return &Storage{db: db}, nil
}

// Сохраняет страницу в хранилище
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return fmt.Errorf("невозможно сохранить страницу: %w", err)
	}

	return nil
}

// Выбирает случайную страницу из хранилища
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("невозможно выбрать случайную страницу: %w", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

// Удаляет страницу из хранилища
func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`
	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return fmt.Errorf("невозможно удалить страницу: %w", err)
	}

	return nil
}

// Проверяет, существует ли страница в хранилище
func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("невозможно проверить существование страницы: %w", err)
	}

	return count > 0, nil
}

// Инициализирует структуру базы данных
func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("невозможно создать таблицу: %w", err)
	}

	return nil
}
