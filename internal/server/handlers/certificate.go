package handlers

import (
	"net/http"

	"github.com/Mario-F/cert-manager-selfservice/internal/gen/api"
	"github.com/Mario-F/cert-manager-selfservice/internal/kube"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (h *OpenAPIV1HandlerImpl) GetCertificateDomain(ctx echo.Context, domain string) error {
	res := &api.Certificate{}
	res.Domain = domain

	certResult, err := kube.GetCertificate(domain, true)
	if err != nil {
		log.Errorf("Error: %v+", err)
		return http.ErrAbortHandler
	}

	// Todo: this businesslogic should be handled by kube/cert
	if len(certResult.CertsFound) == 0 {
		log.Infof("No certs found, creating new cert for domain %s with %s of name %s", domain, IssuerRef.Kind, IssuerRef.Name)
		err := kube.CreateCertificate(domain, IssuerRef)
		if err != nil {
			log.Errorf("Error: %v+", err)
			return http.ErrAbortHandler
		}
		return ctx.NoContent(http.StatusAccepted)
	}

	if len(certResult.CertsFound) > 1 {
		log.Errorf("More than one certificate found: %d", len(certResult.CertsFound))
		return http.ErrAbortHandler
	}

	cert := certResult.CertsFound[0]
	log.Infof("Cert found %s", cert.Certificate.Name)
	if !cert.Ready {
		log.Infof("Cert for domain %s is not ready yet", domain)
		return ctx.NoContent(http.StatusAccepted)
	}

	secretData := cert.Secret.Data
	if len(secretData["tls.key"]) == 0 || len(secretData["tls.crt"]) == 0 {
		log.Errorf("Not all secret data required found for domain: %s", domain)
		return echo.NewHTTPError(http.StatusInternalServerError, "There was a problem fetching certificate secret, see server logs for details")
	}

	res.Authority = string(secretData["ca.crt"])
	res.Crt = string(secretData["tls.crt"])
	res.Key = string(secretData["tls.key"])

	return ctx.JSON(http.StatusOK, res)
}
