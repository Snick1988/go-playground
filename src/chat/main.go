package main

import (
    "log"
    "flag"
    "time"
    "math/rand"
    "net/http"
    "text/template"
    "path/filepath"
    "sync"
)

type templateHandler struct {
    once     sync.Once
    filename string
    templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    t.once.Do(func() {
        t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
    })
    t.templ.Execute(w, r)
}

func main() {
    rand.Seed(time.Now().Unix())

    var addr = flag.String("addr", ":8080", "The addr of the application.")
    flag.Parse()
    room := newRoom()
    
    http.Handle("/", &templateHandler{filename: "chat.html"})
    http.Handle("/room", room)
    go room.run()

    // start the web server
    log.Println("Starting web server on", *addr)
    if err := http.ListenAndServe(*addr, nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}
