
package api

import (
	"github.com/gorilla/mux"
	"github.com/ryomak/go-generate-example/src/controller"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

type Router struct {
	*mux.Router
}

func New() *Router {
	return &Router{
		mux.NewRouter(),
	}
}

func (r *Router) Route(addr string) {
	a := r.PathPrefix("/api").Subrouter()
	r.NotFoundHandler = http.HandlerFunc(controller.NotFound)
	a.Path("user/:id").Queries(
			"key","key",
			"date","date",
	).HandlerFunc(controller.GetUser).Methods("GET")
		a.Path("/").HandlerFunc(controller.Index).Methods("GET")
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(r)
	log.Println("server start :" + addr)
	http.ListenAndServe(":"+addr, n)
}
