package gorillamux

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	"github.com/gorilla/mux"
)

type MuxRouter struct {
	m   *mux.Router
	log infrastructure.LoggerI
}

func New(log infrastructure.LoggerI) *MuxRouter {
	var m = &MuxRouter{
		m:   mux.NewRouter(),
		log: log,
	}
	m.m.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Error 405 - MethodNotAllowed"))
		w.WriteHeader(http.StatusMethodNotAllowed)
		m.log.Println("StatusMethodNotAllowed")
	})
	m.m.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Error 404 - StatusNotFound " + req.URL.Path))
		w.WriteHeader(http.StatusNotFound)
		m.log.Println("NotFoundHandler")
	})
	return m
}

func (r *MuxRouter) Router() http.Handler {
	return r.m
}

func (r *MuxRouter) PathHandler(
	tpl string,
	handler http.Handler,
) infrastructure.RouterI {
	r.m.PathPrefix(tpl).Handler(handler)
	return r
}

func (r *MuxRouter) PathHandlerFunc(
	tpl string,
	f func(http.ResponseWriter, *http.Request),
) infrastructure.RouterI {
	r.m.PathPrefix(tpl).HandlerFunc(f)
	return r
}

func (r *MuxRouter) PathSubrouter(
	tpl string,
) infrastructure.RouterI {
	return &MuxRouter{
		m:   r.m.PathPrefix(tpl).Subrouter(),
		log: r.log,
	}
}

func (r *MuxRouter) GET(
	path string,
	f handlers.ResultFunc,
) infrastructure.RouterI {
	r.m.HandleFunc(path, handlers.HandleFunc(f)).Methods("GET")
	return r
}

func (r *MuxRouter) POST(
	path string,
	f handlers.ResultFunc,
) infrastructure.RouterI {
	r.m.HandleFunc(path, handlers.HandleFunc(f)).Methods("POST")
	return r
}

func (r *MuxRouter) PUT(
	path string,
	f handlers.ResultFunc,
) infrastructure.RouterI {
	r.m.HandleFunc(path, handlers.HandleFunc(f)).Methods("PUT")
	return r
}

func (r *MuxRouter) DELETE(
	path string,
	f handlers.ResultFunc,
) infrastructure.RouterI {
	r.m.HandleFunc(path, handlers.HandleFunc(f)).Methods("DELETE")
	return r
}

func (r *MuxRouter) OPTIONS(
	path string,
	f handlers.ResultFunc,
) infrastructure.RouterI {
	r.m.HandleFunc(path, handlers.HandleFunc(f)).Methods("OPTIONS")
	return r
}

func (r *MuxRouter) AddMiddleware(
	mwf ...infrastructure.MiddlewareI,
) infrastructure.RouterI {
	var muxMwf = make([]mux.MiddlewareFunc, 0)
	for _, m := range mwf {
		muxMwf = append(muxMwf, m.Func)
	}
	r.m.Use(muxMwf...)
	return r
}
