package auth

import (
	"context"
	"errors"

	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/cantylv/service-happy-birthday/internal/repository/user"
	"github.com/cantylv/service-happy-birthday/internal/utils/functions"
	"github.com/cantylv/service-happy-birthday/internal/utils/myerrors"
)

type Usecase interface {
	SignUpUser(ctx context.Context, data entity.SignUpForm) (string, error)
	SignInUser(ctx context.Context, data entity.SignInForm) (string, error)
}

type UsecaseLayer struct {
	repo user.Repo
}

// NewUsecaseLayer
// Returns an instance of usecase layer.
func NewUsecaseLayer(repo user.Repo) UsecaseLayer {
	return UsecaseLayer{
		repo: repo,
	}
}

func (uc *UsecaseLayer) SignUpUser(ctx context.Context, data entity.SignUpForm) (string, error) {

	u, err := uc.repo.GetByEmail(ctx, data.Email)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserNotExist) {
			uId, err := uc.repo.Create(ctx, functions.ConverterCreateUserDB(&data))
			if err != nil {
				return "", err
			}
			return uId, nil
		}
		return "", err
	}
	return u.Id.Hex(), myerrors.ErrUserAlreadyExist
}

func (uc *UsecaseLayer) SignInUser(ctx context.Context, data entity.SignInForm) (string, error) {
	u, err := uc.repo.GetByEmail(ctx, data.Email)
	if err != nil {
		return "", err
	}
	return u.Id.Hex(), nil
}
