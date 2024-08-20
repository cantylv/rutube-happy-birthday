package crontasks

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/cantylv/service-happy-birthday/internal/repository/user"
	"github.com/cantylv/service-happy-birthday/services/kafka/producer"
	"github.com/cantylv/service-happy-birthday/services/mongodb"
	"github.com/mailru/easyjson"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func InitCronTasks() *cron.Cron {
	c := cron.New()

	c.AddFunc("@every 10s", InitClearTask)
	c.AddFunc("@everyday", BirthdayNotification)

	c.Start()

	return c
}

func InitClearTask() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	// init mongodb connection
	mongoClient := mongodb.Init()
	collection := mongoClient.Database("main").Collection("subs")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			logger.Errorf("error disconnecting from MongoDB: %v", err)
		}
	}()
	filter := bson.M{"subs.is_followed": false}
	update := bson.M{
		"$pull": bson.M{
			"subs": bson.M{"is_followed": false},
		},
	}

	// Выполнение обновления
	res, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Infof("%s was deleted %d", time.Now().Format("02.01.2006 15:04:05 UTC-07"), res.ModifiedCount)
}

func BirthdayNotification() {
	logger := zap.Must(zap.NewProduction()).Sugar()

	prod := producer.Init()
	defer func(prod sarama.SyncProducer) {
		err := prod.Close()
		if err != nil {
			logger.Errorf("error close kafka producer: %v", err)
		}
	}(prod)

	// init mongodb connection
	mongoClient := mongodb.Init()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			logger.Errorf("error disconnecting from MongoDB: %v", err)
		}
	}()
	collection := mongoClient.Database("main").Collection("subs")

	var users []user.User
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		logger.Errorf("error while searching documents: %v", err)
		return
	}
	err = cursor.All(ctx, &users)
	if err != nil {
		logger.Errorf("error while extracting documents using cursor: %v", err)
		return
	}

	makeNotification(prod, users)
}

func makeNotification(prod sarama.SyncProducer, users []user.User) {
	logger := zap.Must(zap.NewProduction())
	timeNow := time.Now().Truncate(24 * time.Hour)
	for _, user := range users {
		for _, sub := range user.Subs {
			// need to calculate dates and time interval between them
			if intervalNotification(sub.EmployeeBirthday, sub.Interval).Equal(timeNow) {
				notification := entity.Notification{
					FollowerEmail:    user.Email,
					FollowerId:       user.Id.Hex(),
					EmployeeEmail:    sub.EmployeeEmail,
					EmployeeId:       sub.EmployeeId,
					EmployeeFullName: sub.EmployeeFullName,
					Interval:         5,
				}
				// Сериализация структуры в JSON
				jsonData, err := easyjson.Marshal(notification)
				if err != nil {
					log.Fatalf("Failed to marshal notification: %v", err)
				}

				msg := sarama.ProducerMessage{
					Topic: viper.GetString("kafka.topic"),
					Value: sarama.ByteEncoder(jsonData),
				}
				_, _, err = prod.SendMessage(&msg)
				if err != nil {
					logger.Error(err.Error())
				}
			}
		}
	}
}

func intervalNotification(employeeBirthday string, interval uint16) time.Time {
	t, _ := time.Parse("", employeeBirthday)
	return t.AddDate(0, 0, int(interval))
}
