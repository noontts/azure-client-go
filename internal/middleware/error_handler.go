package middleware

import (
	"azureclient/internal/errs"
	"encoding/json"
	"net/http"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func ErrorHandler(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}
		if appErr, ok := err.(*errs.AppError); ok {
			w.Header().Set("Content-Type", "application/json")
			code := appErr.Code
			if code == 0 {
				code = http.StatusInternalServerError
			}
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(appErr)
			return
		}
		// fallback for unknown errors
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  5000,
			"message": err.Error(),
		})
	}
}
