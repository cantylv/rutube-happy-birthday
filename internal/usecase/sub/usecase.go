package sub

import (
	"context"
	"errors"

	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/cantylv/service-happy-birthday/internal/repository/sub"
	"github.com/cantylv/service-happy-birthday/internal/repository/user"
	"github.com/cantylv/service-happy-birthday/internal/utils/functions"
	"github.com/cantylv/service-happy-birthday/internal/utils/myerrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Usecase interface {
	Subscribe(ctx context.Context, ids entity.SubProps) error
	Unsubscribe(ctx context.Context, ids entity.SubProps) error
	ChangeInterval(ctx context.Context, intervalData entity.SetUpIntervalProps) error
}

type UsecaseLayer struct {
	repoSub  sub.Repo
	repoUser user.Repo
}

func NewUsecaseLayer(rSub sub.Repo, rUser user.Repo) UsecaseLayer {
	return UsecaseLayer{
		repoSub:  rSub,
		repoUser: rUser,
	}
}

func (uc *UsecaseLayer) Subscribe(ctx context.Context, ids entity.SubProps) error {
	// check follower existence
	followerId, err := primitive.ObjectIDFromHex(ids.IdFollower)
	if err != nil {
		return err
	}
	_, err = uc.repoUser.GetById(ctx, followerId)
	if err != nil {
		return err
	}
	// check employee existence
	employeeId, err := primitive.ObjectIDFromHex(ids.IdEmployee)
	if err != nil {
		return err
	}
	uDB, err := uc.repoUser.GetById(ctx, employeeId)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserNotExist) {
			return myerrors.ErrSubscribeNonExistUser
		}
		return err
	}
	// check if IdFollower == IdEmployee
	if ids.IdEmployee == ids.IdFollower {
		return myerrors.ErrSubscribeYourself
	}
	idsDB, err := functions.ConverterIdsDB(ids, uDB)
	if err != nil {
		return err
	}
	res, err := uc.repoSub.UpdateSubscribtion(ctx, idsDB)
	if err != nil {
		return err
	}
	// When follower doesn't have sub on employee.
	if res.MatchedCount == 0 {
		res, err = uc.repoSub.NewSubscription(ctx, idsDB)
		if err != nil {
			return err
		}
		if res.ModifiedCount == 0 {
			return myerrors.ErrUpdateFailed
		}
		return nil
	}
	return nil
}

func (uc *UsecaseLayer) Unsubscribe(ctx context.Context, ids entity.SubProps) error {
	// check follower existence
	followerId, err := primitive.ObjectIDFromHex(ids.IdFollower)
	if err != nil {
		return err
	}
	_, err = uc.repoUser.GetById(ctx, followerId)
	if err != nil {
		return err
	}
	// check employee existence
	employeeId, err := primitive.ObjectIDFromHex(ids.IdEmployee)
	if err != nil {
		return err
	}
	_, err = uc.repoUser.GetById(ctx, employeeId)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserNotExist) {
			return myerrors.ErrUnsubscribeNonExistUser
		}
		return err
	}
	// check if IdFollower == IdEmployee
	if ids.IdEmployee == ids.IdFollower {
		return myerrors.ErrUnsubscribeYourself
	}
	idsDB, err := functions.ConverterIdsDB(ids, &user.User{})
	if err != nil {
		return err
	}
	_, err = uc.repoSub.Unsubscribe(ctx, idsDB)
	if err != nil {
		return err
	}
	// if mongo.UpdateResult.MatchedCount == 0 and no errors --> sub is not exist OK
	// all situation is OK
	return nil
}

func (uc *UsecaseLayer) ChangeInterval(ctx context.Context, intervalData entity.SetUpIntervalProps) error {
	// check follower existence
	followerId, err := primitive.ObjectIDFromHex(intervalData.Ids.IdFollower)
	if err != nil {
		return err
	}
	_, err = uc.repoUser.GetById(ctx, followerId)
	if err != nil {
		return err
	}
	// check employee existence
	employeeId, err := primitive.ObjectIDFromHex(intervalData.Ids.IdEmployee)
	if err != nil {
		return err
	}
	_, err = uc.repoUser.GetById(ctx, employeeId)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserNotExist) {
			return myerrors.ErrSetIntervalNonExistUser
		}
		return err
	}
	// check if IdFollower == IdEmployee
	if intervalData.Ids.IdEmployee == intervalData.Ids.IdFollower {
		return myerrors.ErrSetIntervalYourself
	}
	intervalDataDB, err := functions.ConverterIntervalDB(intervalData)
	if err != nil {
		return err
	}
	// check that follower followed employee
	idsDB, err := functions.ConverterIdsDB(intervalData.Ids, &user.User{})
	if err != nil {
		return err
	}
	isFollow, err := uc.repoSub.IsFollowed(ctx, idsDB)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if !isFollow {
		return myerrors.ErrSetIntervalNotSubscribe
	}
	res, err := uc.repoSub.ChangeInterval(ctx, intervalDataDB)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return myerrors.ErrNoSubscription
	}
	return nil
}
