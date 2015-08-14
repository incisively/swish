package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	adminPort := flag.String("bind", ":8999", "where to bind the admin API. Make sure it's not public.")
	flag.Parse()

	api := NewAPI()
	log.Printf("ADMIN: Starting server on %v...", *adminPort)
	log.Fatal(http.ListenAndServe(*adminPort, api))
}
