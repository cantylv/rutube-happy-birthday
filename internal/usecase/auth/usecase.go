package auth

import (
	"context"
	"errors"

	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/cantylv/service-happy-birthday/internal/repository/user"
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
			uId, err := uc.repo.Create(ctx, &entity.User{
				FullName: data.FullName,
				Email:    data.Email,
				Password: data.Password, // need to hash and salt
				Birthday: data.Birthday, // need to format
			})
			if err != nil {
				return "", err
			}
			return uId, nil
		}
		return "", err
	}
	return u.Id, myerrors.ErrUserAlreadyExist
}

func (uc *UsecaseLayer) SignInUser(ctx context.Context, data entity.SignInForm) (string, error) {
	u, err := uc.repo.GetByEmail(ctx, data.Email)
	if err != nil {
		return "", err
	}
	return u.Id, nil
}
