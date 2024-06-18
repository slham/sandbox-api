package request

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, messages []string) {
	RespondWithJSON(w, code, map[string][]string{"errors": messages})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
