package util

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func SendJSON(w http.ResponseWriter, resp any, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		SendJSON(w, map[string]string{"error": "Internal Server Error"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("Failed to write response", "error", err)
		return
	}
}
