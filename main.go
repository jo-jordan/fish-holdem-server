package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	hub := newHub()

	go hub.run()

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/login", login)
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	log.Fatal(http.ListenAndServe(*addr, nil))
}
