package handlers

import (
	"net/http"

	"github.com/Mario-F/cert-manager-selfservice/internal/config"
	"github.com/Mario-F/cert-manager-selfservice/internal/gen/api"
	"github.com/labstack/echo/v4"
)

func (h *OpenAPIV1HandlerImpl) GetInfo(ctx echo.Context) error {
	res := &api.Info{}
	res.Version = config.Version
	basicAuthActive, _, _ := config.GetBasicAuth()
	res.BasicAuth = basicAuthActive
	return ctx.JSON(http.StatusOK, res)
}
