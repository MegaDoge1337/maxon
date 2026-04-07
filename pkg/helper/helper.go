package helper

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

func ParseID(s string) (int, error) {
	return strconv.Atoi(s)
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode response", slog.Any("error", err))
	}
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}
