package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, code int, object interface{}) {
	var encoded []byte
	var err error
	if object == nil {
		encoded = []byte("{}")
	} else {
		encoded, err = json.Marshal(object)
		if err != nil {
			code = http.StatusInternalServerError
		}
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(encoded)

	log.Printf("Response: %s", encoded)
}
