package database

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
	"tz/models"
)

type List struct {
	*mongo.Collection
	cache
	edited
}

type edited struct {
	sync.RWMutex
	b bool
}

type cache struct {
	sync.RWMutex
	m map[int]*models.Element
}

func NewList(col *mongo.Collection) (*List, error) {
	list := List{
		Collection: col,
		edited:     edited{b: false},
		cache:      cache{m: make(map[int]*models.Element)},
	}

	err := list.Load()
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (db *List) Get() (results map[int]*models.Element, err error) {
	db.edited.Lock()
	defer db.edited.Unlock()

	if db.edited.b {
		err := db.Load()
		if err != nil {
			return nil, err
		}
	}

	db.edited.b = false

	return db.getListFromCache()
}

func (db *List) Edit(el models.Element) (upserted, modified int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := db.ReplaceOne(ctx, bson.D{{Key: "id", Value: el.ID}}, el,
		options.Replace().SetUpsert(true))
	if err != nil {
		return
	}

	upserted = int(res.UpsertedCount)
	modified = int(res.ModifiedCount)

	db.setEdited(true)

	return
}

func (db *List) Remove(id int) (count int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := db.DeleteOne(ctx, bson.D{{"id", id}})
	if err != nil {
		return
	}

	count = int(res.DeletedCount)

	db.setEdited(true)

	return
}

func (db *List) Load() (err error) {
	db.cache.Lock()
	defer db.cache.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//defer cancel()
	//
	//for _, el := range db.m {
	//	find, err := db.Find(ctx, bson.D{{"id", el.ID}})
	//	if err != nil {
	//		return err
	//	}
	//
	//	res := &models.Element{}
	//	err = find.Decode(res)
	//	if err != nil {
	//		return err
	//	}
	//
	//	if res == nil {
	//		delete(db.m, el.ID)
	//	}
	//}

	cursor, err := db.Find(ctx, bson.D{})
	if err != nil {
		return
	}

	res := make([]bson.D, 0)

	err = cursor.All(context.TODO(), &res)
	if err != nil {
		return
	}

	db.cache.m = make(map[int]*models.Element)
	for _, el := range res {
		switch el.Map()["id"].(type) {
		case int64:
			db.m[int(el.Map()["id"].(int64))] = models.FromBson(el)
		case int32:
			db.m[int(el.Map()["id"].(int32))] = models.FromBson(el)
		}

	}

	log.Info().Msg("loaded")

	return err
}

func (c *cache) getListFromCache() (results map[int]*models.Element, err error) {
	c.RLock()
	defer c.RUnlock()

	results = c.m

	return
}

func (e *edited) setEdited(b bool) {
	e.Lock()
	defer e.Unlock()

	e.b = b
}
