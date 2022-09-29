package database

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DB struct {
	*mongo.Database
}

func (db *DB) ListGet() ([]bson.D, error) {
	log.Info().Msg("list get")
	var results []bson.D
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cursor, err := db.Collection("list").Find(ctx, bson.D{})
	if err != nil {
		return results, err
	}

	err = cursor.All(context.TODO(), &results)

	return results, err
}

func (db *DB) ListEdit(doc bson.D) error {
	log.Info().Msg("list edit")
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = db.Collection("list").ReplaceOne(ctx, doc, doc, options.Replace().SetUpsert(true))

	return err
}

func (db *DB) ListRemove(id int) error {
	log.Info().Msg("list remove")
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = db.Collection("list").DeleteOne(ctx, bson.D{{"id", id}})

	return err
}
