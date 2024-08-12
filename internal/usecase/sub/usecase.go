package sub

import (
	"context"
	"errors"

	"github.com/cantylv/service-happy-birthday/internal/repository/sub"
	"github.com/cantylv/service-happy-birthday/internal/utils/myerrors"
)

type Usecase interface {
	Subscribe(context.Context, sub.SubProps) error
	Unsubscribe(context.Context, sub.SubProps) error
	ChangeInterval(context.Context, sub.SetUpIntervalProps) error
}

type UsecaseLayer struct {
	repo sub.Repo
}

func (u *UsecaseLayer) Subscribe(ctx context.Context, ids sub.SubProps) error {
	res, err := u.repo.Subscribe(ctx, ids)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		u.repo.NewSubscription(ctx, ids)
	} else if res.ModifiedCount == 0 {
		return errors.New(myerrors.UpdateFailed)
	}
	return nil
}

func (u *UsecaseLayer) Unsubscribe(ctx context.Context, ids sub.SubProps) error {
	res, err := u.repo.Unsubscribe(ctx, ids)
	if err != nil {
		return err
	}
	if res.MatchedCount != 0 && res.ModifiedCount == 0 {
		return errors.New(myerrors.UpdateFailed)
	}
	return nil
}
