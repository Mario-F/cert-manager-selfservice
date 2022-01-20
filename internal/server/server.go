package server

import (
	"fmt"
	"net/http"

	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

// Start is the entrypoint for starting the webserver
func Start(port int, issuerKind string, issuerName string) {
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

	e.GET("/selfcert/:domain/pem", getSelfCertHandler())

	e.GET("/cert/:domain/:crttype", getCertHandler(issuerRef))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
