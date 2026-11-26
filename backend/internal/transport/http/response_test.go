package http

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

// Record every forwarded status: ResponseRecorder itself suppresses duplicates,
// so using it alone would hide a broken wrapper from this regression test.
type recordingWriter struct {
	header   http.Header
	statuses []int
	body     bytes.Buffer
}

func (w *recordingWriter) Header() http.Header { return w.header }
func (w *recordingWriter) WriteHeader(status int) {
	w.statuses = append(w.statuses, status)
}
func (w *recordingWriter) Write(body []byte) (int, error) {
	return w.body.Write(body)
}

func TestResponseWriterPreservesFirstStatus(t *testing.T) {
	cases := []struct {
		name          string
		initialStatus int
		laterStatus   int
		wantStatus    int
	}{
		{
			name:          "explicit_status_cannot_be_replaced",
			initialStatus: http.StatusCreated,
			laterStatus:   http.StatusInternalServerError,
			wantStatus:    http.StatusCreated,
		},
		{
			name:        "implicit_success_cannot_be_replaced",
			laterStatus: http.StatusBadRequest,
			wantStatus:  http.StatusOK,
		},
		{
			name:          "repeated_status_is_not_forwarded",
			initialStatus: http.StatusOK,
			laterStatus:   http.StatusOK,
			wantStatus:    http.StatusOK,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			underlying := &recordingWriter{header: make(http.Header)}
			writer := &responseWriter{ResponseWriter: underlying}
			if tt.initialStatus != 0 {
				writer.WriteHeader(tt.initialStatus)
			}
			if n, err := writer.Write([]byte("first")); err != nil || n != 5 {
				t.Fatalf("initial write=(%d,%v); want (5,nil)", n, err)
			}
			writer.WriteHeader(tt.laterStatus)
			if n, err := writer.Write([]byte("second")); err != nil || n != 6 {
				t.Fatalf("subsequent write=(%d,%v); want (6,nil)", n, err)
			}
			if !reflect.DeepEqual(underlying.statuses, []int{tt.wantStatus}) {
				t.Fatalf("forwarded statuses=%v; want only [%d]", underlying.statuses, tt.wantStatus)
			}
			if got := underlying.body.String(); got != "firstsecond" {
				t.Fatalf("body=%q; want firstsecond", got)
			}
		})
	}
}
