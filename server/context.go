package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// Context used for network reader interface
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

// NewContext exported function used for network reader interface
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{w, r}
}

// Header exported function used for network reader interface
func (c *Context) Header(name, value string) {
	c.ResponseWriter.Header().Set(name, value)
}

// Param exported function used for network reader interface
func (c *Context) Param(name string) string {
	params := mux.Vars(c.Request)
	return params[name]
}

// RenderError exported function used for network reader interface
func (c *Context) RenderError(status int, err error) {
	http.Error(c.ResponseWriter, err.Error(), status)

}

// SetStatus exported function used for network reader interface
func (c *Context) SetStatus(s int) {
	c.ResponseWriter.WriteHeader(s)
}

// RenderJSON exported function used for network reader interface
func (c *Context) RenderJSON(status int, j interface{}) {
	c.Header("Content-type", "application/json")
	c.SetStatus(status)
	data, err := json.Marshal(j)
	if err != nil {
		c.SetStatus(http.StatusInternalServerError)
		return
	}

	c.ResponseWriter.Write(data)
}
