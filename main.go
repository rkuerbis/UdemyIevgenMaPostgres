package main

import (
	"github.com/udemy_fileserver/api"
	"github.com/udemy_fileserver/datastore"
	"log"
	"runtime"

	"github.com/udemy_fileserver/server"
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
