package errors

import (
	"encoding/json"
	"net/http"
)

// Lets us centralize writing errors in a JSON response for easy handling by callers/scripts
func WriteJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}