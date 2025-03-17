package handlers

import (
	r "blog-api/repositories/blog"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

/*
ParseBlogQueryParams modifies a pointer to BlogQuery in place.
It looks up the required query values from the url.Values map,
parsing those values and loading them into the BlogQuery struct
if present, otherwise using the default initialized values from
the struct.
*/
func ParseBlogQueryParams(q *r.BlogQuery, v url.Values) {
	offset := v.Get("offset")

	if offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err == nil {
			q.Offset = parsedOffset
		}
	}
}

func WriteJSONErr(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if encoderErr := json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()}); encoderErr != nil {
		handleFallBackResponse(w)
	}
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if encoderErr := json.NewEncoder(w).Encode(data); encoderErr != nil {
		handleFallBackResponse(w)
	}
}

func handleFallBackResponse(w http.ResponseWriter) {
	fallback := `{"error": "Internal Server Error - failed to encode response"}`

	if _, writeErr := w.Write([]byte(fallback)); writeErr != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: failed to write response: %v", writeErr), http.StatusInternalServerError)
		return
	}
}

func EmptyResponse() map[string]interface{} {
	emptyResponse := map[string]interface{}{}
	return emptyResponse
}
