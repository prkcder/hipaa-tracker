package handlers


import (
    "fmt"
    "net/http"

)


func RootHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "🚀 HIPAA Tracker API is up")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    // fmt.Fprintln(w, "ok")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
}



