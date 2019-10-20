package main

import (
	"fmt"
	"net/http"
)

const (
	port = "8080"
	host = "http://127.0.0.1:" + port
)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		h.handleMain(w, r)
	case "/login":
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	case "/result":
		h.handleResult(w, r)
	case "/download":
		h.handleDownLoad(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	fields := []string{"BEGIN:", "SUMMARY:", "DTSTART;VALUE=DATE:", "DTEND;VALUE=DATE:", "RRULE:FREQ=YEARLY;UNTIL=", "DESCRIPTION:", "END:"}
	handler := &Handler{
		IcsFields: fields,
	}

	fmt.Println("Start listening at port: ", port)
	http.ListenAndServe(":"+port, handler)
}
