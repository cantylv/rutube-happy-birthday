// Copyright © ivanlobanov. All rights reserved.
package user

import (
	"context"
	"errors"

	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/cantylv/service-happy-birthday/internal/repository/user"
	"github.com/cantylv/service-happy-birthday/internal/utils/functions"
	"github.com/cantylv/service-happy-birthday/internal/utils/myerrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Usecase interface {
	GetData(ctx context.Context, employee_id string) (*entity.User, error)
	UpdateData(ctx context.Context, uData *entity.UserUpdate, uId string) error
}

type UsecaseLayer struct {
	repoUser user.Repo
}

func NewUsecaseLayer(rUser user.Repo) UsecaseLayer {
	return UsecaseLayer{
		repoUser: rUser,
	}
}

func (uc *UsecaseLayer) GetData(ctx context.Context, employee_id string) (*entity.User, error) {
	objectId, err := primitive.ObjectIDFromHex(employee_id)
	if err != nil {
		return nil, err
	}
	uDB, err := uc.repoUser.GetById(ctx, objectId)
	if err != nil {
		return nil, err
	}
	u, err := functions.ConverterUserEntity(uDB)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *UsecaseLayer) UpdateData(ctx context.Context, uData *entity.UserUpdate, uId string) error {
	u, err := uc.repoUser.GetByEmail(ctx, uData.Email)
	if err != nil && !errors.Is(err, myerrors.ErrUserNotExist) {
		return err
	}
	if u != nil && uId != u.Id.Hex() {
		return myerrors.ErrEmailIsReserved
	}
	objectId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return err
	}
	err = uc.repoUser.Update(ctx, functions.ConverterUpdateUserDB(uData, objectId))
	if err != nil {
		return err
	}
	return nil
}
