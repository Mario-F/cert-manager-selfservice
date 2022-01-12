package server

import (
	"net/http"

	"github.com/Mario-F/cert-manager-selfservice/internal/kube"
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type PlainCert struct {
	Ca  string `json:"ca"`
	Crt string `json:"crt"`
	Key string `json:"key"`
}

func getCertHandler(issuerRef cmmeta.ObjectReference, certPrefix string) echo.HandlerFunc {
	return func(c echo.Context) error {
		domain := c.Param("domain")
		crttype := c.Param("crttype")
		log.Infof("Certservice called with domain: %s", domain)

		certResult, err := kube.GetCertificate(domain, true)
		if err != nil {
			log.Errorf("Error: %v+", err)
			return http.ErrAbortHandler
		}

		if len(certResult.CertsFound) == 0 {
			log.Infof("No certs found, creating new cert for domain %s with %s of name %s", domain, issuerRef.Kind, issuerRef.Name)
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
		if crttype == "json" {
			pCrt := PlainCert{
				string(secretData["ca.crt"]),
				string(secretData["tls.crt"]),
				string(secretData["tls.key"]),
			}
			err := c.JSON(200, pCrt)
			if err != nil {
				log.Errorf("Error: %v+", err)
			}
		}

		return nil
	}
}
