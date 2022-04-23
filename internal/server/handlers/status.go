package handlers

import (
	"net/http"
	"time"

	"github.com/Mario-F/cert-manager-selfservice/internal/gen/api"
	"github.com/labstack/echo/v4"
)

func (h *OpenAPIV1HandlerImpl) GetStatus(ctx echo.Context) error {
	res := &api.Status{}
	// TODO: Fill with real data
	res.Messages = append(res.Messages, api.StatusMessage{
		Message:  "This is a placeholder message for now",
		Severity: "info",
		Time:     time.Now(),
	})
	return ctx.JSON(http.StatusOK, res)
}
