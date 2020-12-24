package handlers

import (
	"mathbattle/models/mathbattle"
	"net/http"

	"github.com/gorilla/mux"
)

type ProblemHandler struct {
	Ps mathbattle.ProblemService
}

func (h *ProblemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]

	problem, err := h.Ps.GetByID(ID)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, problem)
}
