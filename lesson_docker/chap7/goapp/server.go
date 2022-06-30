package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "hello"}`))
		default:
			http.Error(w, ``, http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(":8888", nil)
}
