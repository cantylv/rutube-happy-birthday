// Copyright © ivanlobanov. All rights reserved.
package sub

import (
	"context"

	"github.com/cantylv/service-happy-birthday/internal/utils/myconstants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo interface {
	UpdateSubscribtion(ctx context.Context, ids SubProps) (*mongo.UpdateResult, error)
	NewSubscription(ctx context.Context, ids SubProps) (*mongo.UpdateResult, error)
	Unsubscribe(ctx context.Context, ids SubProps) (*mongo.UpdateResult, error)
	ChangeInterval(ctx context.Context, data SetUpIntervalProps) (*mongo.UpdateResult, error)
}

type RepoLayer struct {
	cl *mongo.Collection
}

type SubProps struct {
	IdFollower primitive.ObjectID
	IdEmployee string
}

type SetUpIntervalProps struct {
	Ids         SubProps
	NewInterval uint16
}

// NewRepoLayer
// Returns an instance of repository layer.
func NewRepoLayer(collection *mongo.Collection) RepoLayer {
	return RepoLayer{
		cl: collection,
	}
}

// UpdateSubscribtion
// Update subscribtion to an employee. Result --> element in array 'subs' with field 'is_followed == true'.
func (r *RepoLayer) UpdateSubscribtion(ctx context.Context, ids SubProps) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": ids.IdFollower, "subs.employee_id": ids.IdEmployee}
	newData := bson.M{"$set": bson.M{"subs.$.is_followed": true}}
	return r.cl.UpdateOne(ctx, filter, newData)
}

// NewSubscription
// Subscribes to an employee. Result --> new element in array 'subs'.
func (r *RepoLayer) NewSubscription(ctx context.Context, ids SubProps) (*mongo.UpdateResult, error) {
	newData := bson.M{
		"employee_id": ids.IdEmployee,
		"interval":    myconstants.DefaultInterval,
		"is_followed": true,
	}
	filter := bson.M{"_id": ids.IdFollower}
	newData = bson.M{
		"$push": bson.M{"subs": newData},
	}
	return r.cl.UpdateOne(ctx, filter, newData)
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
func (r *RepoLayer) ChangeInterval(ctx context.Context, data SetUpIntervalProps) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": data.Ids.IdFollower, "subs.employee_id": data.Ids.IdEmployee}
	newData := bson.M{"$set": bson.M{"subs.$.interval": data.NewInterval}}
	return r.cl.UpdateOne(ctx, filter, newData)
}
