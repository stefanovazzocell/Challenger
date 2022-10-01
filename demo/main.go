package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Preparing demo server...")
	handler := NewHandler()
	http.HandleFunc("/", handler.handleDemo)
	http.HandleFunc("/challenge.js", handler.handleJs)
	http.HandleFunc("/challenge", handler.handleChallenge)
	http.HandleFunc("/solve", handler.handleSolve)

	log.Println("Starting demo server at http://localhost:8888")
	if err := http.ListenAndServe("localhost:8888", nil); err != nil {
		panic(err)
	}
}
