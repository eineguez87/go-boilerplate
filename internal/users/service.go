package users

import (
	"context"
	"errors"

	"go-boilerplate/internal/auth"
)

type Service interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
	Me(ctx context.Context, userID string) (User, error)
}

type service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) Service {
	return &service{repo: repo, jwtSecret: jwtSecret}
}

func (s *service) Register(ctx context.Context, email, password string) error {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	return s.repo.Create(ctx, email, hash)
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	id, hash, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := auth.CheckPassword(hash, password); err != nil {
		return "", errors.New("invalid credentials")
	}

	return auth.GenerateToken(id, s.jwtSecret)
}

func (s *service) Me(ctx context.Context, userID string) (User, error) {
	return s.repo.FindByID(ctx, userID)
}
