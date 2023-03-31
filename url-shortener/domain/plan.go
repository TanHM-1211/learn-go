package domain

import (
	"context"
	"time"
)

type Plan struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Desc       string    `json:"description"`
	MaxNum     int       `json:"maximum_number"`
	MaxDuraion time.Time `json:"max_duration"`
}

type PlanRepository interface {
	GetByID(ctx context.Context, id int) (Plan, error)
}
