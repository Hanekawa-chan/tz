package models

import "go.mongodb.org/mongo-driver/bson"

type Element struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func FromBson(d bson.D) *Element {
	id := 0
	switch d.Map()["id"].(type) {
	case int64:
		id = int(d.Map()["id"].(int64))
	case int32:
		id = int(d.Map()["id"].(int32))
	}
	return &Element{
		ID:   id,
		Text: d.Map()["text"].(string),
	}
}
