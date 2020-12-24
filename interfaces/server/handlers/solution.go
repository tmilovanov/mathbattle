package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"mathbattle/models/mathbattle"

	"github.com/gorilla/mux"
)

type SolutionHandler struct {
	Ss mathbattle.SolutionService
}

func (h *SolutionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var solution mathbattle.Solution
	err := json.NewDecoder(r.Body).Decode(&solution)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	solution, err = h.Ss.Create(solution)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, solution)
}

func (h *SolutionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]

	solution, err := h.Ss.Get(ID)
	if err != nil {
		if err == mathbattle.ErrNotFound {
			ResponseJSON(w, http.StatusNotFound, nil)
			return
		}
	}

	ResponseJSON(w, http.StatusOK, solution)
}

func (h *SolutionHandler) Find(w http.ResponseWriter, r *http.Request) {
	var findDescriptor mathbattle.FindDescriptor
	err := json.NewDecoder(r.Body).Decode(&findDescriptor)
	if err != nil {
		log.Printf("Failed to decode findDescriptor, error: %v", err)
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	solutions, err := h.Ss.Find(findDescriptor)
	if err != nil {
		log.Printf("Failed to find solutions by descriptor, error: %v", err)
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, solutions)
}

func (h *SolutionHandler) AppendPart(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]

	var part mathbattle.Image
	err := json.NewDecoder(r.Body).Decode(&part)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	err = h.Ss.AppendPart(ID, part)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, nil)
}

func (h *SolutionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]

	err := h.Ss.Delete(ID)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, nil)
}

func (h *SolutionHandler) GetProblemDescriptors(w http.ResponseWriter, r *http.Request) {
	participantID := mux.Vars(r)["participant_id"]

	desc, err := h.Ss.GetProblemDescriptors(participantID)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, desc)
}
