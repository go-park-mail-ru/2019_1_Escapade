package service

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/handlers"
)

type Router struct {
	h *handlers.GameHandler
}

func (r *Router) Init(h *handlers.GameHandler) *Router {
	r.h = h
	return r
}

func (r *Router) Router() http.Handler {
	return r.h.Router()
}
