package server

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/Mario-F/cert-manager-selfservice/internal/logger"
)

// Start is the entrypoint for starting the webserver
func Start(port int) {
	logger.Verbosef("Starting webserver...")

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		logger.Verbosef("default handler called")
		fmt.Fprintf(rw, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
