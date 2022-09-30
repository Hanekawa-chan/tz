package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"tz/models"
)

type DB struct {
	*mongo.Database
}

func (db *DB) ListGet() (results []bson.D, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cursor, err := db.Collection("list").Find(ctx, bson.D{})
	if err != nil {
		return
	}

	err = cursor.All(context.TODO(), &results)

	return
}

func (db *DB) ListEdit(el models.Element) (upserted, modified int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := db.Collection("list").
		ReplaceOne(ctx, bson.D{{Key: "id", Value: el.ID}}, el,
			options.Replace().SetUpsert(true))
	if err != nil {
		return
	}

	upserted = int(res.UpsertedCount)
	modified = int(res.ModifiedCount)

	return
}

func (db *DB) ListRemove(id int) (count int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := db.Collection("list").DeleteOne(ctx, bson.D{{"id", id}})
	if err != nil {
		return
	}

	count = int(res.DeletedCount)

	return
}
