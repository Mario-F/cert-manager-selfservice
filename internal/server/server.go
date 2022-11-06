package server

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/Mario-F/cert-manager-selfservice/internal/gen/api"
	"github.com/Mario-F/cert-manager-selfservice/internal/kube"
	"github.com/Mario-F/cert-manager-selfservice/internal/server/handlers"
	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
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
func Start(port int, issuerKind string, issuerName string, authUsername string, authPassword string) {
	log.Infof("Starting webserver with IssuerKind: %s and IssuerName: %s", issuerKind, issuerName)
	issuerRef := cmmeta.ObjectReference{
		Name: issuerName,
		Kind: issuerKind,
	}

	e = echo.New()
	e.Use(echoPrometheus.MetricsMiddleware())

	e.Use(middleware.Logger())

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "web/dist",
		Index:      "index.html",
		Browse:     false,
		HTML5:      true,
		Filesystem: http.FS(EmbededStatic),
	}))

	e.GET("/selfcert/:domain/pem", getSelfCertHandler())

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Errorf("Error loading swagger spec\n: %s", err)
	}
	e.GET("/api/spec/v1", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, swagger)
	})
	OpenapiHandlerImpl := &handlers.OpenAPIV1HandlerImpl{}
	validatorOptions := &oapimiddleware.Options{}
	validatorOptions.Options.AuthenticationFunc = func(c context.Context, input *openapi3filter.AuthenticationInput) error {
		log.Infof("Authenticating called with %s but not implemented at the moment", input.SecuritySchemeName)

		// Validate if authUsername and authPassword is set
		if input.SecuritySchemeName == "basicAuth" && (authUsername != "" || authPassword != "") {
			username, password, ok := input.RequestValidationInput.Request.BasicAuth()
			if !ok || username != authUsername || password != authPassword {
				log.Errorf("Auth failed for user %s", username)
				return fmt.Errorf("invalid credentials")
			}
			log.Debugf("Auth successful for user %s", username)
			return nil
		}

		// If no validation method activated all credentials work
		return nil
	}
	apiGroup := e.Group("/api/v1", oapimiddleware.OapiRequestValidatorWithOptions(swagger, validatorOptions))
	kube.IssuerRef = issuerRef
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
