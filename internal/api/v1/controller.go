package v1

import (
	"encoding/json"
	"net/http"
)

type Controller struct{}

func (Controller) GetPing(w http.ResponseWriter, r *http.Request) {
	response := Pong{
		Ping: "pong",
	}

	respondWithJSON(w, http.StatusOK, response)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
