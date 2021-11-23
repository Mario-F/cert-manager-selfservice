package server

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/Mario-F/cert-manager-selfservice/internal/cert"
	"github.com/Mario-F/cert-manager-selfservice/internal/logger"
	"github.com/gorilla/mux"
)

// Start is the entrypoint for starting the webserver
func Start(port int) {
	logger.Verbosef("Starting webserver...")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		logger.Verbosef("default handler called")
		//logger.Debugf("Request %v+", r)
		vars := mux.Vars(r)
		logger.Debugf("Request %v+", vars)
		fmt.Fprintf(rw, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	myRouter.HandleFunc("/cert/{domain}", func(rw http.ResponseWriter, r *http.Request) {
		logger.Verbosef("cert handler called")
		//logger.Debugf("Request %v+", r)
		vars := mux.Vars(r)
		logger.Debugf("Request %v+", vars)
		fmt.Fprintf(rw, "Request cert: %s\n", vars["domain"])
		certRequest := cert.SelfSignedCertRequest{
			Domain: vars["domain"],
		}
		cert, err := cert.SelfSignedCert(certRequest)
		if err != nil {
			logger.Errorf("Error: %v+", err)
		}
		fmt.Fprintf(rw, "Public: %s\n", &cert.CertPEM)
		fmt.Fprintf(rw, "Private: %s\n", &cert.PrivatePEM)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), myRouter))
}
