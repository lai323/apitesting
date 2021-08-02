package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithField("Method", r.Method).WithField("Path", r.URL.Path).Info("request")
		next.ServeHTTP(w, r)
	})
}
