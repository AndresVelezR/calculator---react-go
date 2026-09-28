package httpapi

import "net/http"

func CORS(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Vary", "Origin")
		origin := request.Header.Get("Origin")
		if origin != "" && origin == allowedOrigin {
			writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}
		if request.Method == http.MethodOptions && origin != "" {
			if origin != allowedOrigin {
				writeJSON(writer, http.StatusForbidden, ErrorResponse{Error: "origin is not allowed"})
				return
			}
			method := request.Header.Get("Access-Control-Request-Method")
			if method != http.MethodPost && method != http.MethodGet {
				writeJSON(writer, http.StatusMethodNotAllowed, ErrorResponse{Error: "method is not allowed"})
				return
			}
			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
