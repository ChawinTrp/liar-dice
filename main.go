package main

import (
	"log"
	"net/http"
)

func main() {
	// 1. Register the WebSocket handler
	// When a client requests /ws, serveWs will upgrade the connection
	http.HandleFunc("/ws", serveWs)

	// 2. Optional: A simple health check endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Liar's Dice Game Server is Running!"))
	})

	// 3. Start the HTTP server on port 8080
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
