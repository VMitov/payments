package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8000", "address:port")
	db := flag.String("db", "postgres://postgres@localhost:5432", "address:port")
	flag.Parse()

	api, err := newAPI(*db)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(*addr, newRouter(api)))
}
