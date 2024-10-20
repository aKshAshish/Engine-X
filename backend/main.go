package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.Handle("/", http.HandlerFunc(baseHandler))
	log.Println("Starting server on port:", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal("error starting the server", err)
	}
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("PORT")
	res := map[string]string{
		"msg": fmt.Sprintf("Hello from %s!", port),
	}
	w.Header().Set("conten-type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("error sending the response")
	}
}
