package service

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/handlers"
)

type Router struct {
	h *handlers.Handlers
}

func (r *Router) Init(h *handlers.Handlers) *Router {
	r.h = h
	return r
}

func (r *Router) Router() http.Handler {
	return r.h.Router()
}
