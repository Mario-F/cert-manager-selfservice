package handlers

import (
	"net/http"

	"github.com/Mario-F/cert-manager-selfservice/internal/gen/api"
	"github.com/labstack/echo/v4"
)

func (h *OpenAPIV1HandlerImpl) GetCertificateDomain(ctx echo.Context, domain string) error {
	res := &api.Certificate{}

	return ctx.JSON(http.StatusOK, res)
}
