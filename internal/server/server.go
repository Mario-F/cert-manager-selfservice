package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Mario-F/cert-manager-selfservice/internal/static"
	echoPrometheus "github.com/globocom/echo-prometheus"
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

var e *echo.Echo

// Start is the entrypoint for starting the webserver
func Start(port int, issuerKind string, issuerName string) {
	log.Infof("Starting webserver with IssuerKind: %s and IssuerName: %s", issuerKind, issuerName)
	issuerRef := cmmeta.ObjectReference{
		Name: issuerName,
		Kind: issuerKind,
	}

	e = echo.New()
	e.Use(echoPrometheus.MetricsMiddleware())

	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		log.Debug("default handler called")
		return c.HTML(http.StatusOK, string(static.Get("/default.html")))
	})

	e.GET("/selfcert/:domain/pem", getSelfCertHandler())

	e.GET("/cert/:domain/:crttype", getCertHandler(issuerRef))

	if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}
}

func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := e.Shutdown(ctx)
	if err != nil {
		log.Errorf("Error occured on echo shutdown %v+", err)
	}
}
