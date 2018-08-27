package api

import (
	"github.com/UdemyIevgenMaPostgres/server"
)

// Start exported function capitalized
func Start(r *server.Router) {
	r.OPTIONS("/{rest:.*", server.ConstructRequest)
	handleFiles(r.Group("/files"))
	r.PutToAllRoutes(server.ConstructRequest)
}
