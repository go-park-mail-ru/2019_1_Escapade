package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/middleware"
)

// Use set handlers of http.StatusMethodNotAllowed and
// http.StatusNotFound errors. Also it add Recovery middleware
// and all middleware specified as the parameters
// IMPORTANT: use only with root router!
func Use(r *mux.Router, mwf ...mux.MiddlewareFunc) {
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Error 405 - MethodNotAllowed"))
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("StatusMethodNotAllowed")
	})
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Error 404 - StatusNotFound "+ req.URL.Path))
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("NotFoundHandler")
	})

	r.Use(middleware.Recover)
	r.Use(mwf...)
}
