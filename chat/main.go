package main

import (
	"flag"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
)

type temlateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// handle http for template
func (t *temlateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// will only load this template once and cache it
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The address of the application server")
	flag.Parse()
	r := newRoom()
	http.Handle("/", &temlateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// start the room
	go r.run()

	// start the server
	log.Infof("Listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
