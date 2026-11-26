package http

import (
	"net/http"
	"net/url"
	"strings"
)

func allowedOriginSet(allowedOrigins []string) map[string]bool {
	// Copy configuration once; the immutable map is safe for concurrent requests.
	origins := make(map[string]bool, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		u, err := url.Parse(origin)
		if err == nil && (u.Scheme == "http" || u.Scheme == "https") &&
			u.Host != "" && u.User == nil && u.Path == "" && u.RawQuery == "" && u.Fragment == "" {
			origins[origin] = true
		}
	}
	return origins
}

// handleCORS returns true only when the actual resource handler should run.
// CORS is a browser policy, not authentication: origin-less clients still work.
func handleCORS(w http.ResponseWriter, r *http.Request, origins map[string]bool, methods string) bool {
	w.Header().Add("Vary", "Origin")
	preflight := r.Method == http.MethodOptions
	if preflight {
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")
	}
	origin := r.Header.Get("Origin")
	deny := func() bool {
		writeProblem(w, http.StatusForbidden, "Cross-origin request is not allowed.", codeCORSDenied)
		return false
	}
	if (origin != "" && !origins[origin]) || (preflight && origin == "") {
		return deny()
	}
	if preflight {
		requestMethod := r.Header.Get("Access-Control-Request-Method")
		if !methodAllowed(requestMethod, methods) {
			return deny()
		}
		requestHeaders := r.Header.Get("Access-Control-Request-Headers")
		if requestHeaders != "" {
			for _, header := range strings.Split(requestHeaders, ",") {
				if !strings.EqualFold(strings.TrimSpace(header), "Content-Type") {
					return deny()
				}
			}
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", methods)
		w.WriteHeader(http.StatusNoContent)
		return false
	}
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	return true
}

func methodAllowed(method, allowed string) bool {
	for _, candidate := range strings.Split(allowed, ", ") {
		if method == candidate {
			return true
		}
	}
	return false
}
