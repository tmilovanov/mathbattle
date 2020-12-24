package handlers

import (
	"encoding/json"
	"mathbattle/models/mathbattle"
	"net/http"

	"github.com/gorilla/mux"
)

type ReviewHandler struct {
	Rs mathbattle.ReviewService
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	var review mathbattle.Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	review, err = h.Rs.Store(review)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, review)
}

func (h *ReviewHandler) FindMany(w http.ResponseWriter, r *http.Request) {
	var findDescriptor mathbattle.ReviewFindDescriptor
	err := json.NewDecoder(r.Body).Decode(&findDescriptor)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	reviews, err := h.Rs.FindMany(findDescriptor)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, reviews)
}

func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]

	err := h.Rs.Delete(ID)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, nil)
}

func (h *ReviewHandler) GetSolutionDescriptors(w http.ResponseWriter, r *http.Request) {
	participantID := mux.Vars(r)["participant_id"]

	descs, err := h.Rs.RevewStageDescriptors(participantID)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, descs)
}
