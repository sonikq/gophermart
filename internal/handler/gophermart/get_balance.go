package gophermart

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sonikq/gophermart/internal/models"
	"net/http"
)

func (h *Handler) GetBalance(ctx *gin.Context) {
	const source = "handler.GetBalance"

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

	var response *models.Balance
	response, err = h.service.GetBalance(c, username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error, something went wrong"})
		h.logger.Error().
			Err(err).
			Str("source", source).
			Msg("failed to get balance")
		return
	}

	ctx.JSON(http.StatusOK, response)
}
