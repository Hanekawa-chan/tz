package models

import "go.mongodb.org/mongo-driver/bson"

type Element struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func (e *Element) ToBson() bson.D {
	var d bson.D
	d = bson.D{
		{"id", e.ID},
		{"text", e.Text},
	}
	return d
}
