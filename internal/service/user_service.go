package service

import (
	"context"
	"errors"
	"time"
	"todo_api/internal/domain"
	"todo_api/pkg/hasher"
	"todo_api/pkg/jwtutil"
)

type UserRepository interface {
	Create(ctx context.Context, input domain.CreateUserInput) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userService struct {
	repo      UserRepository
	jwtSecret string
	jwtTTL    time.Duration
}

func NewUserService(repo UserRepository, jwtSecret string, jwtTTL time.Duration) *userService {
	return &userService{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtTTL:    jwtTTL,
	}

}

func (s *userService) Register(ctx context.Context, input domain.CreateUserInput) (*domain.User, error) {
	if len(input.Password) < 6 {
		return nil, domain.ErrInvalidInput
	}
	hashed, err := hasher.Hash(input.Password)

	if err != nil {
		return nil, err
	}

	user := domain.CreateUserInput{
		Email:    input.Email,
		Password: hashed,
	}
	return s.repo.Create(ctx, user)
}

func (s *userService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.ErrNotFound
		}
		return "", err
	}

	if err = hasher.Compare(user.Password, password); err != nil {
		return "", domain.ErrUnauthorized
	}

	token, err := jwtutil.GenerateToken(s.jwtSecret, jwtutil.Claims{
		UserID: user.ID,
		Email:  user.Email,
	}, s.jwtTTL)
	if err != nil {
		return "", err
	}
	return token, nil
}
