package utils

import (
	"encoding/json"
	"net/http"
)

func SendJSON(w http.ResponseWriter, rawData any) {
	data, err := json.Marshal(rawData)
	if err != nil {
		http.Error(w, "failed to serialize json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
