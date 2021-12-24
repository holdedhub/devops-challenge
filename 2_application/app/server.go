package main


import (
    "fmt"
    "log"
    "net/http"
)

type StatusRecorder struct {
    http.ResponseWriter
    Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
    r.Status = status
    r.ResponseWriter.WriteHeader(status)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
    }
    name := r.FormValue("name")
    address := r.FormValue("address")
    fmt.Fprintf(w, "Name = %s\n", name)
    fmt.Fprintf(w, "Address = %s\n", address)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    fmt.Fprintf(w, "Hello!")
}

func logRequest(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        recorder := &StatusRecorder{
            ResponseWriter: w,
            Status: 200,
        }
        h.ServeHTTP(recorder, r)
        log.Println(r.RemoteAddr, r.Host, r.Method, r.RequestURI, r.ContentLength, recorder.Status)
    })
}

func main() {
    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/", fileServer)
    http.HandleFunc("/form", formHandler)
    http.HandleFunc("/hello", helloHandler)


    fmt.Printf("Starting server at port 8080\n")
    if err := http.ListenAndServe(":8080", logRequest(http.DefaultServeMux)); err != nil {
        log.Fatal(err)
    }
}

