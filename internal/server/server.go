package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/Mario-F/cert-manager-selfservice/internal/gen/api"
	"github.com/Mario-F/cert-manager-selfservice/internal/server/handlers"
	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	echoPrometheus "github.com/globocom/echo-prometheus"
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

var (
	e             *echo.Echo
	EmbededStatic *embed.FS
)

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

	staticRoot, err := fs.Sub(EmbededStatic, "web/dist")
	if err != nil {
		e.Logger.Fatal("Error as descent in static subdirectory", err)
	}
	staticHandler := http.FileServer(http.FS(staticRoot))
	e.GET("/*", echo.WrapHandler(staticHandler))

	e.GET("/selfcert/:domain/pem", getSelfCertHandler())

	e.GET("/cert/:domain/:crttype", getCertHandler(issuerRef))

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Errorf("Error loading swagger spec\n: %s", err)
	}
	e.GET("/api/spec/v1", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, swagger)
	})
	OpenapiHandlerImpl := &handlers.OpenAPIV1HandlerImpl{}
	apiGroup := e.Group("/api/v1", oapimiddleware.OapiRequestValidator(swagger))
	api.RegisterHandlers(apiGroup, OpenapiHandlerImpl)

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
