package functions

import (
	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/cantylv/service-happy-birthday/internal/repository/sub"
	"github.com/cantylv/service-happy-birthday/internal/repository/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConverterUpdateUserDB(data *entity.UserUpdate, uId primitive.ObjectID) *user.User {
	return &user.User{
		Id:       uId,
		FullName: data.FullName,
		Birthday: data.Birthday,
		Email:    data.Email,
		Password: data.Password,
	}
}

func ConverterCreateUserDB(data *entity.SignUpForm) *user.User {
	return &user.User{
		FullName: data.FullName,
		Birthday: data.Birthday,
		Email:    data.Email,
		Password: data.Password,
		Subs:     []entity.Subscription{},
	}
}

func ConverterIdsDB(ids entity.SubProps) (sub.SubProps, error) {
	followerId, err := primitive.ObjectIDFromHex(ids.IdFollower)
	if err != nil {
		return sub.SubProps{}, err
	}
	return sub.SubProps{
		IdFollower: followerId,
		IdEmployee: ids.IdEmployee,
	}, nil
}

func ConverterIntervalDB(ids entity.SetUpIntervalProps) (sub.SetUpIntervalProps, error) {
	followerId, err := primitive.ObjectIDFromHex(ids.Ids.IdFollower)
	if err != nil {
		return sub.SetUpIntervalProps{}, err
	}
	return sub.SetUpIntervalProps{
		Ids: sub.SubProps{
			IdFollower: followerId,
			IdEmployee: ids.Ids.IdEmployee,
		},
		NewInterval: ids.NewInterval,
	}, nil
}

func ConverterUserEntity(data *user.User) *entity.User {
	return &entity.User{
		Id:       data.Id.Hex(),
		FullName: data.FullName,
		Birthday: data.Birthday,
		Email:    data.Email,
		Password: data.Password,
		Subs:     data.Subs,
	}
}
