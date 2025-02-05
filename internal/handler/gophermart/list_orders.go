package gophermart

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sonikq/gophermart/internal/models"
	"net/http"
)

func (h *Handler) ListOrders(ctx *gin.Context) {
	const source = "handler.ListOrders"

	username, err := getUsername(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": models.ErrNotAuthenticated.Error()})
		h.logger.Error().
			Err(err).
			Str("source", source).
			Send()
		return
	}

	c, cancel := context.WithTimeout(ctx, h.config.CtxTimeOut)
	defer cancel()

	var response []models.Order
	response, err = h.service.ListUserOrders(c, username)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": "Internal server error, something went wrong"})
		h.logger.Error().Err(err).Str("source", source).Msg("failed to list orders")
		return
	}

	if len(response) == 0 {
		ctx.AbortWithStatusJSON(
			http.StatusNoContent,
			gin.H{"error": "Empty result"})
		h.logger.Error().Str("source", source).Msg("empty response")
		return
	}

	ctx.JSON(http.StatusOK, response)
}
