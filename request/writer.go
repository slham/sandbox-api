package request

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"errors": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	slog.Debug("RESPONSE", "status", code, "headers", w.Header())
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
