// Copyright © ivanlobanov. All rights reserved.
package sub

import (
	"context"
	"errors"

	"github.com/cantylv/service-happy-birthday/internal/utils/myconstants"
	"github.com/cantylv/service-happy-birthday/internal/utils/myerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo interface {
	Subscribe(ctx context.Context, ids SubProps) (*mongo.UpdateResult, error)
	NewSubscription(ctx context.Context, ids SubProps) error
	Unsubscribe(ctx context.Context, ids SubProps) (*mongo.UpdateResult, error)
	ChangeInterval(ctx context.Context, data SetUpIntervalProps) error
}

type RepoLayer struct {
	cl *mongo.Collection
}

// NewRepoLayer
// Returns an instance of repository layer.
func NewRepoLayer(collection *mongo.Collection) RepoLayer {
	return RepoLayer{
		cl: collection,
	}
}

// Props
type SubProps struct {
	IdFollower uint32
	IdEmployee uint32
}

type SetUpIntervalProps struct {
	Ids         SubProps
	NewInterval uint16
}

// Subscribe
// Subscribes to an employee. Result --> element in array 'subs' with field 'is_followed == true'.
func (r *RepoLayer) Subscribe(ctx context.Context, ids SubProps) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": ids.IdFollower, "subs.employee_id": ids.IdEmployee}
	newData := bson.M{"$set": bson.M{"subs.$.is_followed": true}}
	return r.cl.UpdateOne(ctx, filter, newData)
}

// NewSubscription
// Subscribes to an employee. Result --> new element in array 'subs'.
func (r *RepoLayer) NewSubscription(ctx context.Context, ids SubProps) error {
	newData := bson.M{
		"employee_id": ids.IdEmployee,
		"interval":    myconstants.DefaultInterval,
		"is_followed": true,
	}
	filter := bson.M{"_id": ids.IdFollower}
	newData = bson.M{
		"$push": bson.M{"subs": newData},
	}
	resUpdate, err := r.cl.UpdateOne(ctx, filter, newData)
	if err != nil {
		return err
	}
	if resUpdate.MatchedCount == 0 || resUpdate.ModifiedCount == 0 {
		return errors.New(myerrors.UpdateFailed)
	}
	return nil
}

// Unsubscribe
// Unsubscribes to an employee. Result --> element in array has field `is_followed` with value `false`.
// Cron task will remove all records with `is_followed==false`.
func (r *RepoLayer) Unsubscribe(ctx context.Context, ids SubProps) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": ids.IdFollower, "subs.employee_id": ids.IdEmployee}
	newData := bson.M{"$set": bson.M{"subs.$.is_followed": false}}
	return r.cl.UpdateOne(ctx, filter, newData)
}

// ChangeInterval
// Change the value of field 'subs.interval' for specific element in field (array) 'subs'.
func (r *RepoLayer) ChangeInterval(ctx context.Context, data SetUpIntervalProps) error {
	filter := bson.M{"_id": data.Ids.IdFollower, "subs.employee_id": data.Ids.IdEmployee}
	newData := bson.M{"$set": bson.M{"subs.$.interval": data.NewInterval}}
	res, err := r.cl.UpdateOne(ctx, filter, newData)
	if err != nil {
		return err
	}
	if res.MatchedCount != 0 && res.ModifiedCount == 0 {
		return errors.New(myerrors.UpdateFailed)
	}
	return nil
}
