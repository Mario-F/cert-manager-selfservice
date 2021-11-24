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

	myRouter.HandleFunc("/cert/{domain}/pem", func(rw http.ResponseWriter, r *http.Request) {
		logger.Verbosef("cert handler called")
		vars := mux.Vars(r)
		logger.Debugf("Request %v+", vars)
		certRequest := cert.SelfSignedCertRequest{
			Domain: vars["domain"],
		}
		cert, err := cert.SelfSignedCert(certRequest)
		if err != nil {
			logger.Errorf("Error: %v+", err)
		}

		// Output certificate
		_, err = rw.Write(cert.PrivatePEM.Bytes())
		if err != nil {
			logger.Errorf("Error: %v+", err)
		}
		_, err = rw.Write(cert.CertPEM.Bytes())
		if err != nil {
			logger.Errorf("Error: %v+", err)
		}
		_, err = rw.Write(cert.CaCertPEM.Bytes())
		if err != nil {
			logger.Errorf("Error: %v+", err)
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), myRouter))
}
