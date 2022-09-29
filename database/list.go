package database

import "go.mongodb.org/mongo-driver/mongo"

type DB struct {
	*mongo.Database
}

func (db *DB) ListGet() *mongo.Collection {
	col := db.Collection("list")
	return col
}
