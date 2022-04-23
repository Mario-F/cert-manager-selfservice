package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (h *OpenAPIV1HandlerImpl) GetStatus(ctx echo.Context) error {
	log.Debug("GetStatus")
	return ctx.JSON(http.StatusOK, "OK")
}
