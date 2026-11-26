package http

import (
	"log/slog"
	"net/http"

	"github.com/martinGM05/sezzle-technical-assessment/backend/internal/calculator"
)

const (
	healthPath       = "/health"
	calculationsPath = "/api/v1/calculations"
)

// NewHandler returns the calculator HTTP boundary. Cross-origin requests are
// disabled unless their exact HTTP(S) origin is explicitly supplied. No credentials
// or wildcard origins are supported; invalid configured origins are ignored.
func NewHandler(allowedOrigins ...string) http.Handler {
	origins := allowedOriginSet(allowedOrigins)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := &responseWriter{ResponseWriter: w}
		w = response
		if r.Method == http.MethodHead {
			w = headWriter{ResponseWriter: w}
		}
		defer func() {
			if recovered := recover(); recovered != nil {
				if recovered == http.ErrAbortHandler || response.started {
					panic(http.ErrAbortHandler)
				}
				slog.Error("request handler panicked")
				writeProblem(w, http.StatusInternalServerError, "An unexpected error occurred.", codeInternalError)
			}
		}()

		var allow string
		switch r.URL.Path {
		case healthPath:
			allow = "GET, HEAD"
		case calculationsPath:
			allow = "POST"
		default:
			writeProblem(w, http.StatusNotFound, "Route does not exist.", codeNotFound)
			return
		}

		if !handleCORS(w, r, origins, allow) {
			return
		}
		if !methodAllowed(r.Method, allow) {
			w.Header().Set("Allow", allow)
			writeProblem(w, http.StatusMethodNotAllowed, "Method is not supported for this route.", codeMethodNotAllowed)
			return
		}

		if r.URL.Path == healthPath {
			w.Header().Set("Content-Type", mediaTypeJSON)
			writeJSON(w, healthResponse{Status: "ok"})
			return
		}
		calculate(w, r)
	})
}

func calculate(w http.ResponseWriter, r *http.Request) {
	operation, operands, valid := decodeRequest(w, r)
	if !valid {
		return
	}

	result, calculationError := calculator.Calculate(operation, operands)
	if calculationError != nil {
		writeProblem(w, http.StatusUnprocessableEntity, calculationError.Detail, calculationError.Code)
		return
	}
	w.Header().Set("Content-Type", mediaTypeJSON)
	writeJSON(w, calculationResponse{Result: result})
}
