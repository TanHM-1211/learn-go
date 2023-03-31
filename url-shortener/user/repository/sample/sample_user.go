package sample

import (
	"context"

	"github.com/TanHM-1211/url_shortener/domain"
)

type sampleUserRepository struct{}

func NewSampleUserRepository() domain.UserRepository {
	return &sampleUserRepository{}
}

func (s *sampleUserRepository) GetByID(ctx context.Context, cursor string) (domain.User, error) {
	return domain.User{}, nil
}
