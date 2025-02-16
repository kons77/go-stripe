package main

import "net/http"

func SessioLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
