package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
	"tz/models"
)

func (h *Handler) ListGet(w http.ResponseWriter, r *http.Request) {
	var err error

	list, err := h.DB.List.Get()
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

	upserted, modified, err := h.DB.List.Edit(element)
	if err != nil {
		log.Error().Err(err).Msg("db edit")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]int{"upserted": upserted, "modified": modified})
	if err != nil {
		log.Error().Err(err).Msg("json encode")
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

	count, err := h.DB.List.Remove(id)
	if err != nil {
		log.Error().Err(err).Msg("db edit")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]int{"count": count})
	if err != nil {
		log.Error().Err(err).Msg("json encode")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
