package domain

import (
	"context"
	"time"
)

type URLRecord struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"original_url"`
	UserID      string    `json:"user_id"`
	Expiration  time.Time `json:"expiration"`
}

type URLRecordUseCase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]URLRecord, string, error)
	GetByID(ctx context.Context, id string) (URLRecord, error)
	Update(ctx context.Context, url *URLRecord) error
	Store(ctx context.Context, url *URLRecord) error
	Delete(ctx context.Context, id string) error
}

type URLRecordRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]URLRecord, string, error)
	GetByID(ctx context.Context, id string) (URLRecord, error)
	Update(ctx context.Context, url *URLRecord) error
	Store(ctx context.Context, url *URLRecord) error
	Delete(ctx context.Context, id string) error
}
