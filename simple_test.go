package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting simple HTTP server...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test endpoint working!")
	})

	fmt.Println("Server starting on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
