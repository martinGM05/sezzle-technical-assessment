package http

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"github.com/martinGM05/sezzle-technical-assessment/backend/internal/calculator"
)

const maxRequestBodyBytes = 1 << 20

type request struct {
	Operation *string     `json:"operation"`
	Operands  *[]*float64 `json:"operands"`
}

func decodeRequest(w http.ResponseWriter, r *http.Request) (string, []float64, bool) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != mediaTypeJSON {
		writeProblem(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json.", codeUnsupportedMediaType)
		return "", nil, false
	}

	body, readErr := io.ReadAll(http.MaxBytesReader(w, r.Body, maxRequestBodyBytes))
	if readErr != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(readErr, &sizeErr) {
			writeProblem(w, http.StatusRequestEntityTooLarge, "Request body exceeds 1048576 bytes.", codeBodyTooLarge)
		} else {
			writeProblem(w, http.StatusBadRequest, "Request body could not be read.", codeMalformedJSON)
		}
		return "", nil, false
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil || fields == nil {
		writeProblem(w, http.StatusBadRequest, "Request body must be valid JSON.", codeMalformedJSON)
		return "", nil, false
	}
	for key := range fields {
		if key != "operation" && key != "operands" {
			writeProblem(w, http.StatusBadRequest, "Request contains an unknown property.", codeMalformedJSON)
			return "", nil, false
		}
	}

	var req request
	if raw, ok := fields["operation"]; ok {
		if err := json.Unmarshal(raw, &req.Operation); err != nil {
			writeProblem(w, http.StatusBadRequest, "Operation must be a string.", codeMalformedJSON)
			return "", nil, false
		}
	}
	if raw, ok := fields["operands"]; ok {
		if err := json.Unmarshal(raw, &req.Operands); err != nil {
			writeProblem(w, http.StatusBadRequest, "Operands must be an array of numbers.", codeMalformedJSON)
			return "", nil, false
		}
	}
	if req.Operation == nil || req.Operands == nil {
		writeProblem(w, http.StatusUnprocessableEntity, "Operation and operands are required.", calculator.CodeInvalidOperands)
		return "", nil, false
	}

	operands := make([]float64, len(*req.Operands))
	for i, v := range *req.Operands {
		if v == nil {
			writeProblem(w, http.StatusUnprocessableEntity, "Operands must not be null.", calculator.CodeInvalidOperands)
			return "", nil, false
		}
		operands[i] = *v
	}
	return *req.Operation, operands, true
}
