package auth

import (
	"context"

	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/cantylv/service-happy-birthday/internal/repository/user"
)

type Usecase interface {
	SignUp(data *entity.User) error
	SignIn(data *entity.User) error
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

func (u *UsecaseLayer) SignUp(ctx context.Context, data *entity.User) {
	// проверка авторизован пользователь или нет
	// проверка на то, есть ли такой пользователь
	// если есть, то возвращаем ошибку
	// создаем пользователя
	// проставляем куки
}
