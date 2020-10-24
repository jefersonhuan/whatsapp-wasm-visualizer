package main

import (
	"log"
	"net/http"
)

func main() {
	log.Printf("listening on 8080...")
	err := http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))
	log.Fatalln(err)
}
