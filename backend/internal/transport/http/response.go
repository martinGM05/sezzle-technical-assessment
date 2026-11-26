package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	mediaTypeJSON            = "application/json"
	mediaTypeProblem         = "application/problem+json"
	codeInternalError        = "INTERNAL_ERROR"
	codeNotFound             = "NOT_FOUND"
	codeMethodNotAllowed     = "METHOD_NOT_ALLOWED"
	codeUnsupportedMediaType = "UNSUPPORTED_MEDIA_TYPE"
	codeBodyTooLarge         = "BODY_TOO_LARGE"
	codeMalformedJSON        = "MALFORMED_JSON"
	codeCORSDenied           = "CORS_DENIED"
)

type problem struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
	Code   string `json:"code"`
}

type healthResponse struct {
	Status string `json:"status"`
}

type calculationResponse struct {
	Result float64 `json:"result"`
}

type headWriter struct {
	http.ResponseWriter
}

func (w headWriter) Write(body []byte) (int, error) {
	return len(body), nil
}

func writeProblem(w http.ResponseWriter, status int, detail, code string) {
	title := http.StatusText(status)
	if status == http.StatusUnprocessableEntity {
		title = "Unprocessable Content"
	}
	w.Header().Set("Content-Type", mediaTypeProblem)
	w.WriteHeader(status)
	writeJSON(w, problem{
		Type:   "about:blank",
		Title:  title,
		Status: status,
		Detail: detail,
		Code:   code,
	})
}

// Once headers are sent, recovery must abort rather than append a second JSON
// document. The standard server owns connection and request-body cleanup.
type responseWriter struct {
	http.ResponseWriter
	started bool
}

func (w *responseWriter) WriteHeader(status int) {
	if w.started {
		return
	}
	w.started = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(body []byte) (int, error) {
	if !w.started {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(body)
}

func writeJSON(w http.ResponseWriter, value any) {
	if err := json.NewEncoder(w).Encode(value); err != nil {
		// A failed write may already have sent bytes. Retrying with a Problem
		// Details response would corrupt the response rather than repair it.
		slog.Error("response write failed", "error", err)
	}
}
