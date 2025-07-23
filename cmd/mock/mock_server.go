// cmd/mock_server.go
// used to test fowarding
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/mock-endpoint", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		log.Printf("📩 Mock API received request:\n%s", string(body))
		fmt.Println(string(body))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Received")
	})

	log.Println("🚀 Mock API listening on http://localhost:9000/mock-endpoint")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
