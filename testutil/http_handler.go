package testutil

import (
	"net/http"
)

func NewTestHandler(c int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(c)
	})
}
