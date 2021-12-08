package server

import (
	"fmt"
	"net/http"

	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/Mario-F/cert-manager-selfservice/internal/cert"
	"github.com/Mario-F/cert-manager-selfservice/internal/kube"
)

// Start is the entrypoint for starting the webserver
func Start(port int, certPrefix string, issuerKind string, issuerName string) {
	log.Infof("Starting webserver with IssuerKind: %s and IssuerName: %s", issuerKind, issuerName)
	issuerRef := cmmeta.ObjectReference{
		Name: issuerName,
		Kind: issuerKind,
	}

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		log.Debug("default handler called")
		return c.String(http.StatusOK, "default route")
	})

	e.GET("/selfcert/:domain/pem", func(c echo.Context) error {
		log.Infof("selfcert handler called")
		domain := c.Param("domain")
		log.Debugf("Request domain %s", domain)

		certRequest := cert.SelfSignedCertRequest{
			Domain: domain,
		}
		cert, err := cert.SelfSignedCert(certRequest)
		if err != nil {
			log.Errorf("Error: %v+", err)
		}

		// Output certificate
		rw := c.Response().Writer
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

		return nil
	})

	e.GET("/cert/:domain/:crttype", func(c echo.Context) error {
		domain := c.Param("domain")
		crttype := c.Param("crttype")
		log.Infof("Certservice called with domain: %s", domain)

		certResult, err := kube.GetCertificate(domain, true)
		if err != nil {
			log.Errorf("Error: %v+", err)
			return http.ErrAbortHandler
		}

		if len(certResult.CertsFound) == 0 {
			log.Infof("No certs found, creating new cert for domain %s with %s of name %s", domain, issuerKind, issuerName)
			err := kube.CreateCertificate(domain, issuerRef, certPrefix)
			if err != nil {
				log.Errorf("Error: %v+", err)
				return http.ErrAbortHandler
			}
			return c.NoContent(http.StatusAccepted)
		}

		if len(certResult.CertsFound) > 1 {
			log.Errorf("More than one certificate found: %d", len(certResult.CertsFound))
			return http.ErrAbortHandler
		}

		cert := certResult.CertsFound[0]
		log.Infof("Cert found %s", cert.Certificate.Name)
		if !cert.Ready {
			log.Infof("Cert for domain %s is not ready yet", domain)
			return c.NoContent(http.StatusAccepted)
		}

		secretData := cert.Secret.Data
		if len(secretData["tls.key"]) == 0 || len(secretData["tls.crt"]) == 0 {
			log.Errorf("Not all secret data required found for domain: %s", domain)
			return echo.NewHTTPError(http.StatusInternalServerError, "There was a problem fetching certificate secret, see server logs for details")
		}

		rw := c.Response().Writer
		outputKey := func() {
			_, err = rw.Write(secretData["tls.key"])
			if err != nil {
				log.Errorf("Error: %v+", err)
			}
		}
		outputCrt := func() {
			_, err = rw.Write(secretData["tls.crt"])
			if err != nil {
				log.Errorf("Error: %v+", err)
			}
		}
		outputCa := func() {
			if len(secretData["ca.crt"]) > 0 {
				_, err = rw.Write(secretData["ca.crt"])
				if err != nil {
					log.Errorf("Error: %v+", err)
				}
			}
		}

		if crttype == "pem" || crttype == "key" {
			outputKey()
		}
		if crttype == "pem" || crttype == "crt" {
			outputCrt()
		}
		if crttype == "pem" || crttype == "ca" {
			outputCa()
		}

		return nil
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
