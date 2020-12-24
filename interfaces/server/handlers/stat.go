package handlers

import (
	"net/http"

	"mathbattle/models/mathbattle"
)

type StatHandler struct {
	Ss mathbattle.StatService
}

func (h *StatHandler) Stat(w http.ResponseWriter, r *http.Request) {
	stat, err := h.Ss.Stat()
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
	}

	ResponseJSON(w, http.StatusOK, stat)
}
