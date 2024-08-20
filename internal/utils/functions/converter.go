package functions

import (
	"encoding/hex"

	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/cantylv/service-happy-birthday/internal/repository/sub"
	"github.com/cantylv/service-happy-birthday/internal/repository/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConverterUpdateUserDB(data *entity.UserUpdate, uId primitive.ObjectID) *user.User {
	// need to hash and salt password, email
	pwdEncoded := hex.EncodeToString([]byte(data.Password))
	emailEncoded := hex.EncodeToString([]byte(data.Email))
	return &user.User{
		Id:       uId,
		FullName: data.FullName,
		Birthday: data.Birthday,
		Email:    emailEncoded,
		Password: pwdEncoded,
	}
}

func ConverterCreateUserDB(data *entity.SignUpForm) *user.User {
	// need to hash and salt password, email
	pwdEncoded := hex.EncodeToString([]byte(data.Password))
	emailEncoded := hex.EncodeToString([]byte(data.Email))
	return &user.User{
		FullName: data.FullName,
		Birthday: data.Birthday,
		Email:    emailEncoded,
		Password: pwdEncoded,
		Subs:     []entity.Subscription{},
	}
}

func ConverterIdsDB(ids entity.SubProps, employeeData *user.User) (sub.SubProps, error) {
	followerId, err := primitive.ObjectIDFromHex(ids.IdFollower)
	if err != nil {
		return sub.SubProps{}, err
	}
	return sub.SubProps{
		IdFollower:       followerId,
		IdEmployee:       ids.IdEmployee,
		FullNameEmployee: employeeData.Email,
		BirthdayEmployee: employeeData.Birthday,
		EmailEmployee:    employeeData.Email,
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

func ConverterUserEntity(data *user.User) (*entity.User, error) {
	decodedEmail, err := hex.DecodeString(data.Email)
	if err != nil {
		return nil, err
	}
	return &entity.User{
		Id:       data.Id.Hex(),
		FullName: data.FullName,
		Birthday: data.Birthday,
		Email:    string(decodedEmail),
		Subs:     data.Subs,
	}, nil
}

func ConverterUserWithoutPwd(data *entity.User) *entity.UserWithoutPassword {
	return &entity.UserWithoutPassword{
		Id:       data.Id,
		FullName: data.FullName,
		Birthday: data.Birthday,
		Email:    data.Email,
		Subs:     data.Subs,
	}
}
