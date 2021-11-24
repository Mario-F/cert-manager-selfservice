package server

import (
	"fmt"
	"html"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Mario-F/cert-manager-selfservice/internal/cert"
	"github.com/gorilla/mux"
)

// Start is the entrypoint for starting the webserver
func Start(port int) {
	log.Info("Starting webserver...")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Debug("default handler called")
		vars := mux.Vars(r)
		log.Debugf("Request %v+", vars)
		fmt.Fprintf(rw, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	myRouter.HandleFunc("/cert/{domain}/pem", func(rw http.ResponseWriter, r *http.Request) {
		log.Infof("cert handler called")
		vars := mux.Vars(r)
		log.Debugf("Request %v+", vars)
		certRequest := cert.SelfSignedCertRequest{
			Domain: vars["domain"],
		}
		cert, err := cert.SelfSignedCert(certRequest)
		if err != nil {
			log.Errorf("Error: %v+", err)
		}

		// Output certificate
		_, err = rw.Write(cert.PrivatePEM.Bytes())
		if err != nil {
			log.Errorf("Error: %v+", err)
		}
		_, err = rw.Write(cert.CertPEM.Bytes())
		if err != nil {
			log.Errorf("Error: %v+", err)
		}
		_, err = rw.Write(cert.CaCertPEM.Bytes())
		if err != nil {
			log.Errorf("Error: %v+", err)
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), myRouter))
}
