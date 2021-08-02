package middleware

import (
	"apitesting/handler"
	"apitesting/token"
	"fmt"
	"net/http"
)

var TokenCodec = token.NewTokenCodec([]byte("AQdJz7ZGbv8GepLS"), []byte{'_'})

func UserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := r.Header.Get("token")
		if t == "" {
			handler.WriteJsonresp(handler.JsonrespInterForbidden(fmt.Errorf("need login")), w)
			return
		}

		_, err := TokenCodec.Token([]byte(t))
		if err != nil {
			handler.WriteJsonresp(handler.JsonrespInterForbidden(err), w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenstr := r.Header.Get("token")
		t, err := TokenCodec.Token([]byte(tokenstr))
		if err != nil {
			handler.WriteJsonresp(handler.JsonrespInterForbidden(err), w)
			return
		}
		if t.Name != "admin" {
			handler.WriteJsonresp(handler.JsonrespInterForbidden(fmt.Errorf("not admin")), w)
		}

		next.ServeHTTP(w, r)
	})
}
