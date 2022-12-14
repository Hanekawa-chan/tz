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
	var cred options.Credential
	time.Sleep(2 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cred.Username = os.Getenv("MONGO_ROOT_USERNAME")
	cred.Password = os.Getenv("MONGO_ROOT_PASSWORD")
	log.Info().Msg(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetAuth(cred))
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

	db := client.Database("tz")
	list, err := database.NewList(db.Collection("list"))
	if err != nil {
		log.Fatal().Err(err).Msg("db list")
	}

	h.DB = &database.DB{List: list}
	go h.Update()

	router := mux.NewRouter()

	h.Router = router
	h.Routes()

	log.Info().Msg("serving on :8080")

	err = http.ListenAndServe(":8080", h.Router)
	if err != nil {
		log.Fatal().Err(err).Msg("listen")
	}
}

func (h *Handler) Routes() {
	h.Router.HandleFunc("/list/get", h.ListGet)
	h.Router.HandleFunc("/list/edit", h.ListEdit)
	h.Router.HandleFunc("/list/remove/{id}", h.ListRemove)
}

func (h *Handler) Update() {
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			go func() {
				log.Info().Msg("trying to load")
				_, err := h.DB.List.Get()
				if err != nil {
					log.Fatal().Err(err).Msg("list cache load")
				}
			}()
		}
	}
}
