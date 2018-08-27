package main

import (
	"github.com/UdemyIevgenMaPostgres/api"
	"github.com/UdemyIevgenMaPostgres/datastore"
	"log"
	"runtime"

	"github.com/UdemyIevgenMaPostgres/server"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)

	datastore.ConnectPostgre()

	s := server.NewRouter()
	api.Start(s.Group("/api"))
	s.Static("/", "./assets/")
	log.Println("Started on :8080")

	if err := s.ListenAndServe("127.0.0.1:8080"); err != nil {
		panic(err)
	}
}
