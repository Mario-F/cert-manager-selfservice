package server

import (
	"fmt"
	"net/http"

	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/Mario-F/cert-manager-selfservice/internal/cert"
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

	e.GET("/cert/:domain/:crttype", getCertHandler(issuerRef, certPrefix))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
