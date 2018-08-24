package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

var allRoutes []func(*Context)

// Router struct type exportable
type Router struct {
	router *mux.Router
}

// NewRouter struct type exportable
func NewRouter() *Router {
	return &Router{router: mux.NewRouter()}
}

// ListenAndServe function exportable
func (r *Router) ListenAndServe(addr string) error {
	srv := &http.Server{
		Handler: r.router,
		Addr:    addr,
	}
	return srv.ListenAndServe()
}

// Group function exportable
func (r *Router) Group(tpl string) *Router {
	return &Router{router: r.router.PathPrefix(tpl).Subrouter()}
}


func (r *Router) route(path, method string, handler func(*Context)) {
	h := func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, r)
		for _, ar := range allRoutes {
			ar(c)

		}
		handler(c)
	}
	if path == "" {
		r.router.Methods(method).HandlerFunc(h)
	} else {
		r.router.Methods(method).Path(path).HandlerFunc(h)
	}

}

// GET network function 
func (r *Router) GET(path string, handler func(*Context)) {
	r.route(path, "GET", handler)
}

// POST network function 
func (r *Router) POST(path string, handler func(*Context)) {
	r.route(path, "POST", handler)
}

// OPTIONS network function 
func (r *Router) OPTIONS(path string, handler func(*Context)) {
	r.route(path, "OPTIONS", handler)
}

// PutToAllRoutes network function 
func (r *Router) PutToAllRoutes(f func(*Context)) {
	allRoutes = append(allRoutes, f)
}

// ConstructRequest network request function
func ConstructRequest(c *Context) {
	c.Header("Access-Control-Allow-Methods", "GET, POST")

	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Max-Age", "86400")

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

}

// DestructRequest network function
func DestructRequest(c *Context) {
	c.Header("Access-Control-Allow-Origin", "*")
}

// Static network function
func (r *Router) Static(pathPrefix, path string) {
	r.router.PathPrefix(pathPrefix).Handler(http.FileServer(http.Dir(path)))
}
