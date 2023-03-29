package server

import "net/http"

type halloHandler struct {
}

func (h *halloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hallo world"))
}

func HalloHandler() http.Handler {
	return &halloHandler{}
}
