package handlers

import (
	"net/http"

	"github.com/Mario-F/cert-manager-selfservice/internal/gen/api"
	"github.com/Mario-F/cert-manager-selfservice/internal/kube"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (h *OpenAPIV1HandlerImpl) GetCertificateDomain(ctx echo.Context, domain string, params api.GetCertificateDomainParams) error {
	res := &api.Certificate{}
	res.Domain = domain

	certResult, err := kube.GetCertificate(domain, true, true)
	if err != nil {
		log.Errorf("Error: %v+", err)
		return http.ErrAbortHandler
	}

	if len(certResult.CertsFound) > 1 {
		log.Errorf("More than one certificate found: %d", len(certResult.CertsFound))
		return http.ErrAbortHandler
	}

	if len(certResult.CertsFound) == 0 {
		if certResult.Created {
			log.Infof("Cert for domain %s has not existst but was created now", domain)
			return ctx.NoContent(http.StatusAccepted)
		}
		log.Infof("Cert for domain %s does not exists and cannot be created", domain)
		return ctx.NoContent(http.StatusNotFound)
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

	log.Debugf("Output format is: %s", *params.Format)
	switch *params.Format {
	case "crt":
		return ctx.HTML(http.StatusOK, res.Crt)
	case "key":
		return ctx.HTML(http.StatusOK, res.Key)
	case "ca":
		return ctx.HTML(http.StatusOK, res.Authority)
	default:
		return ctx.JSON(http.StatusOK, res)
	}
}
