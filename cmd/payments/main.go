package main

import (
	"flag"
	"log"
)

func main() {
	addr := flag.String("addr", ":8000", "address:port")
	flag.Parse()
	log.Fatal(api(*addr))
}
