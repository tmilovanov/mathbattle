package handlers

import (
	"encoding/json"
	"mathbattle/models/mathbattle"
	"net/http"
)

type PostmanHandler struct {
	Ps mathbattle.PostmanService
}

func (h *PostmanHandler) SendToUsers(w http.ResponseWriter, r *http.Request) {
	var msg mathbattle.SimpleMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	err = h.Ps.SendSimpleToUsers(msg)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	ResponseJSON(w, http.StatusOK, nil)
}
