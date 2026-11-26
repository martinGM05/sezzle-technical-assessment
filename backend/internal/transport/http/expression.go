package http

import (
	"encoding/json"
	"net/http"

	"github.com/martinGM05/sezzle-technical-assessment/backend/internal/expression"
)

func evaluateExpression(w http.ResponseWriter, r *http.Request) {
	fields, valid := decodeObject(w, r, "expression")
	if !valid {
		return
	}
	var input *string
	if raw, exists := fields["expression"]; exists {
		if err := json.Unmarshal(raw, &input); err != nil {
			writeProblem(w, http.StatusBadRequest, "Expression must be a string.", codeMalformedJSON)
			return
		}
	}
	if input == nil {
		writeProblem(w, http.StatusUnprocessableEntity, "Expression is required.", expression.CodeInvalidExpression)
		return
	}
	result, err := expression.Evaluate(*input)
	if err != nil {
		writeProblem(w, http.StatusUnprocessableEntity, err.Detail, err.Code)
		return
	}
	w.Header().Set("Content-Type", mediaTypeJSON)
	writeJSON(w, calculationResponse{Result: result})
}
