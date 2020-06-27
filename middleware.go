package main

import (
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
	"time"
)

func (ws *WebServer) nosurfing(h http.Handler) http.Handler {
	surfing := nosurf.New(h)
	surfing.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	return surfing
}

func (ws *WebServer) logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s %s\n", time.Now().Format("02.01.2006 15:04:05"), r.Method, r.URL.Path, r.Proto)
		h.ServeHTTP(w, r)
	})
}

func (ws *WebServer) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/login" || r.URL.Path == "/logout" {
			next.ServeHTTP(w, r)
		} else {
			userInter, err := ws.ab.LoadCurrentUser(&r)
			if userInter != nil && err == nil {
				if userInter.(*User).Role == ROLE_ADMIN {
					next.ServeHTTP(w, r)
				} else {
					if r.URL.Path == "/foo" || r.URL.Path == "/bar" {
						next.ServeHTTP(w, r)
					} else {
						http.Error(w, "Forbidden", http.StatusForbidden)
					}
				}
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		}
	})
}

