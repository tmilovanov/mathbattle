package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"mathbattle/models/mathbattle"

	"github.com/gorilla/mux"
)

type RoundHandler struct {
	Rs mathbattle.RoundService
}

func (h *RoundHandler) StartNew(w http.ResponseWriter, r *http.Request) {
	var startOrder mathbattle.StartOrder
	err := json.NewDecoder(r.Body).Decode(&startOrder)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	round, err := h.Rs.StartNew(startOrder)
	if err != nil {
		log.Printf("Failed to start new round, error: '%v'", err)
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, round)
}

func (h *RoundHandler) GetReivewStageDistribution(w http.ResponseWriter, r *http.Request) {
	desc, err := h.Rs.ReviewStageDistributionDesc()
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, desc)
}

func (h *RoundHandler) StartReviewStage(w http.ResponseWriter, r *http.Request) {
	var startOrder mathbattle.StartOrder
	err := json.NewDecoder(r.Body).Decode(&startOrder)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	round, err := h.Rs.StartReviewStage(startOrder)
	if err != nil {
		log.Printf("Failed to start review stage, error: '%v'", err)
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, round)
}

func (h *RoundHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Print("Handler: GetAll")

	rounds, err := h.Rs.GetAll()
	if err != nil {
		log.Printf("Failed to get all rounds, error: '%v'", err)
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, rounds)
}

func (h *RoundHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	log.Print("Handler: GetByID")

	ID := mux.Vars(r)["id"]

	round, err := h.Rs.GetByID(ID)
	if err != nil {
		if err == mathbattle.ErrNotFound {
			ResponseJSON(w, http.StatusNotFound, nil)
		} else {
			log.Printf("Failed to get round by ID='%s', error: '%v'", ID, err)
			ResponseJSON(w, http.StatusInternalServerError, nil)
		}
		return
	}

	ResponseJSON(w, http.StatusOK, round)
}

func (h *RoundHandler) GetRunning(w http.ResponseWriter, r *http.Request) {
	log.Print("Handler: GetRunning")

	round, err := h.Rs.GetRunning()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			ResponseJSON(w, http.StatusNotFound, nil)
		} else {
			log.Printf("Failed to get current round, error: '%v'", err)
			ResponseJSON(w, http.StatusInternalServerError, nil)
		}
		return
	}

	ResponseJSON(w, http.StatusOK, round)
}

func (h *RoundHandler) GetReviewPending(w http.ResponseWriter, r *http.Request) {
	log.Print("Handler: GetReviewPending")

	round, err := h.Rs.GetReviewPending()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			ResponseJSON(w, http.StatusNotFound, nil)
		} else {
			log.Printf("Failed to get review pending round, error: '%v'", err)
			ResponseJSON(w, http.StatusInternalServerError, nil)
		}
		return
	}

	ResponseJSON(w, http.StatusOK, round)
}

func (h *RoundHandler) GetReviewRunning(w http.ResponseWriter, r *http.Request) {
	log.Print("Handler: GetReviewRunning")

	round, err := h.Rs.GetReviewRunning()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			ResponseJSON(w, http.StatusNotFound, nil)
		} else {
			log.Printf("Failed to get review running round, error: '%v'", err)
			ResponseJSON(w, http.StatusInternalServerError, nil)
		}
		return
	}

	ResponseJSON(w, http.StatusOK, round)
}

func (h *RoundHandler) GetLast(w http.ResponseWriter, r *http.Request) {
	log.Print("Handler: GetLast")

	round, err := h.Rs.GetLast()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			ResponseJSON(w, http.StatusNotFound, nil)
		} else {
			log.Printf("Failed to get last, error: '%v'", err)
			ResponseJSON(w, http.StatusInternalServerError, nil)
		}
		return
	}

	ResponseJSON(w, http.StatusOK, round)
}

func (h *RoundHandler) GetProblemDescriptors(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler: GetProblemDescriptors")

	participantID := mux.Vars(r)["participant_id"]

	desc, err := h.Rs.GetProblemDescriptors(participantID)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, desc)
}
