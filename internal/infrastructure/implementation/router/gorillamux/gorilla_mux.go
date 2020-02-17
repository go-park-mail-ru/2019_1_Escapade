package gorillamux

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/base/handler"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/gorilla/mux"
)

type MuxRouter struct {
	*handler.Handler

	m     *mux.Router
	log   infrastructure.Logger
	trace infrastructure.ErrorTrace
}

// New instance of MuxRouter
func New(
	logger infrastructure.Logger,
	trace infrastructure.ErrorTrace,
) *MuxRouter {
	//overriding the nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}

	//overriding the nil value of ErrorTrace
	if trace == nil {
		trace = new(infrastructure.ErrorTraceNil)
	}

	var m = &MuxRouter{
		Handler: handler.New(logger, trace),

		m:     mux.NewRouter(),
		log:   logger,
		trace: trace,
	}
	m.m.MethodNotAllowedHandler = http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("Error 405 - MethodNotAllowed"))
			m.log.Println("StatusMethodNotAllowed")
		})
	m.m.NotFoundHandler = http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("Error 404 - StatusNotFound " + req.URL.Path))
			m.log.Println("Path not found:", req.URL.Path)
		})
	return m
}

func (r *MuxRouter) Router() http.Handler {
	return r.m
}

func (r *MuxRouter) PathHandler(
	tpl string,
	handler http.Handler,
) infrastructure.Router {
	r.m.PathPrefix(tpl).Handler(handler)
	return r
}

func (r *MuxRouter) PathHandlerFunc(
	tpl string,
	f func(http.ResponseWriter, *http.Request),
) infrastructure.Router {
	r.m.PathPrefix(tpl).HandlerFunc(f)
	return r
}

func (r *MuxRouter) PathSubrouter(
	tpl string,
) infrastructure.Router {
	var m = New(r.log, r.trace)
	m.m = r.m.PathPrefix(tpl).Subrouter()
	return m
}

func (r *MuxRouter) ANY(
	path string,
	f models.ResultFunc,
) infrastructure.Router {
	r.m.HandleFunc(path, r.HandleFunc(f))
	return r
}

func (r *MuxRouter) GET(
	path string,
	f models.ResultFunc,
) infrastructure.Router {
	r.m.HandleFunc(path, r.HandleFunc(f)).Methods("GET")
	return r
}

func (r *MuxRouter) POST(
	path string,
	f models.ResultFunc,
) infrastructure.Router {
	r.m.HandleFunc(path, r.HandleFunc(f)).Methods("POST")
	return r
}

func (r *MuxRouter) PUT(
	path string,
	f models.ResultFunc,
) infrastructure.Router {
	r.m.HandleFunc(path, r.HandleFunc(f)).Methods("PUT")
	return r
}

func (r *MuxRouter) DELETE(
	path string,
	f models.ResultFunc,
) infrastructure.Router {
	r.m.HandleFunc(path, r.HandleFunc(f)).Methods("DELETE")
	return r
}

func (r *MuxRouter) OPTIONS(
	path string,
	f models.ResultFunc,
) infrastructure.Router {
	r.m.HandleFunc(path, r.HandleFunc(f)).Methods("OPTIONS")
	return r
}

func (r *MuxRouter) AddMiddleware(
	mwf ...infrastructure.Middleware,
) infrastructure.Router {
	var muxMwf = make([]mux.MiddlewareFunc, 0)
	for _, m := range mwf {
		muxMwf = append(muxMwf, m.Func)
	}
	r.m.Use(muxMwf...)
	return r
}
