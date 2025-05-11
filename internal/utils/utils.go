package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Envelope is a generic type used for wrapping JSON responses.
// It allows for a consistent structure in API responses, e.g., {"data": ...} or {"error": ...}.
type Envelope map[string]interface{}

// WriteJSON marshals the given data into a JSON string and writes it to the http.ResponseWriter.
// It sets the "Content-Type" header to "application/json" and writes the provided HTTP status code.
// The JSON output is indented for readability.
func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	// Marshal the data to JSON with indentation for pretty printing.
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	// Append a newline character for better formatting in terminals/logs.
	js = append(js, '\n')

	// Set the Content-Type header.
	w.Header().Set("Content-Type", "application/json")

	// Write the HTTP status code.
	w.WriteHeader(status)

	// Write the JSON response body.
	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

// ReadIDParam extracts an "id" URL parameter from the request using Chi's URLParam function.
// It converts the ID to an integer. Returns an error if the parameter is missing or not a valid integer.
func ReadIDParam(r *http.Request) (int, error) {
	// Retrieve the "id" parameter from the URL path.
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		// The parameter was not present.
		return 0, errors.New("invalid id parameter")
	}

	// Convert the string parameter to an integer.
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// The parameter was present but not a valid integer.
		return 0, errors.New("invalid id parameter type")
	}

	// Return the parsed ID.
	return id, nil
}
