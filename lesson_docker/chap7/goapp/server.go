package main

import (
	"log"
	"net/http"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "hello"}`))
	default:
		http.Error(w, ``, http.StatusMethodNotAllowed)
	}
}

func main() {
	port := "8888"

	http.HandleFunc("/", MyHandler)

	done := make(chan bool)
	go http.ListenAndServe(":"+port, nil)
	log.Printf("Server started at port %v", port)
	<-done
}
