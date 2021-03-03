package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	log.Println("server started on [::]:8080")
	if err := http.ListenAndServe(":8080", http.HandlerFunc(handler)); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	e := json.NewEncoder(w)

	res := &greeting{
		Greeting: "Hello, Picard",
	}

	if err := e.Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type greeting struct {
	Greeting string `json:"greeting"`
}
