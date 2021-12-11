package server

import (
	"github.com/Mario-F/cert-manager-selfservice/internal/cert"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func getSelfCertHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
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
	}
}
