package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net/http"
	"strconv"
	"tz/models"
)

func (h *Handler) ListGet(w http.ResponseWriter, r *http.Request) {
	var err error

	list, err := h.DB.ListGet()
	if err != nil {
		log.Error().Err(err).Msg("list get")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		log.Error().Err(err).Msg("json encode")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListEdit(w http.ResponseWriter, r *http.Request) {
	var err error
	var doc bson.D
	var element models.Element

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("read from r.Body")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(b, &element)
	if err != nil {
		log.Error().Err(err).Msg("json unmarshal")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	doc = element.ToBson()

	err = h.DB.ListEdit(doc)
	if err != nil {
		log.Error().Err(err).Msg("db edit")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListRemove(w http.ResponseWriter, r *http.Request) {
	var err error

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Error().Err(err).Msg("id isn't int")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.DB.ListRemove(id)
	if err != nil {
		log.Error().Err(err).Msg("db edit")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
