package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"mathbattle/models/mathbattle"

	"github.com/gorilla/mux"
)

type ParticipantHandler struct {
	Ps mathbattle.ParticipantService
}

func (h *ParticipantHandler) Store(w http.ResponseWriter, r *http.Request) {
	var participant mathbattle.Participant
	err := json.NewDecoder(r.Body).Decode(&participant)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	participant, err = h.Ps.Store(participant)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, participant)
}

func (h *ParticipantHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]

	participant, err := h.Ps.GetByID(ID)
	if err != nil {
		if err == mathbattle.ErrNotFound {
			ResponseJSON(w, http.StatusNotFound, nil)
		} else {
			log.Printf("Failed to get participant by ID='%s', error: '%v'", ID, err)
			ResponseJSON(w, http.StatusInternalServerError, nil)
		}
		return
	}

	ResponseJSON(w, http.StatusOK, participant)
}

func (h *ParticipantHandler) GetByTelegramID(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := mux.Vars(r)["id"]

	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		ResponseJSON(w, http.StatusBadRequest, nil)
		return
	}

	participant, err := h.Ps.GetByTelegramID(telegramID)
	if err != nil {
		if err == mathbattle.ErrNotFound {
			ResponseJSON(w, http.StatusNotFound, nil)
		} else {
			log.Printf("Failed to get participant by telegram ID='%d', error: '%v'", telegramID, err)
			ResponseJSON(w, http.StatusInternalServerError, nil)
		}
		return
	}

	ResponseJSON(w, http.StatusOK, participant)

}

func (h *ParticipantHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	participants, err := h.Ps.GetAll()
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, participants)
}

func (h *ParticipantHandler) Update(w http.ResponseWriter, r *http.Request) {
	ResponseJSON(w, http.StatusInternalServerError, nil)
}

func (h *ParticipantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]

	err := h.Ps.Delete(ID)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, nil)
}

func (h *ParticipantHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]

	err := h.Ps.Delete(ID)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, nil)
}
