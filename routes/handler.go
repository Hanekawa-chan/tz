package routes

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"time"
	"tz/database"
)

type Handler struct {
	Router *mux.Router
	DB     *database.DB
}

func (h *Handler) Run() {
	db := initDB()
	router := mux.NewRouter()

	h.Router = router
	h.DB = &database.DB{Database: db}

	log.Info().Msg("serving on :8080")

	err := http.ListenAndServe("localhost:8080", h.Router)
	if err != nil {
		log.Fatal().Err(err).Msg("listen")
	}
}

func initDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info().Msg(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal().Err(err).Msg("db connect")
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal().Err(err).Msg("db disconnect")
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("db ping")
	}

	return client.Database("tz")
}

func (h *Handler) Routes() {
	h.Router.HandleFunc("/list/get", h.ListGet)
	h.Router.HandleFunc("/list/edit/{id}", h.ListEdit)
	h.Router.HandleFunc("/list/remove/{id}", h.ListRemove)
}
