package gophermart

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sonikq/gophermart/internal/models"
	"github.com/sonikq/gophermart/pkg/reader"
	"net/http"
)

func (h *Handler) Withdraw(ctx *gin.Context) {
	const source = "handler.Withdraw"

	username, err := getUsername(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": models.ErrNotAuthenticated.Error()})
		h.logger.Error().
			Err(err).
			Str("source", source).
			Send()
		return
	}

	if ctx.GetHeader(contentTypeHeaderKey) != contentTypeJSON {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid type of content"})
		h.logger.Error().
			Str("error", "invalid content type").
			Str("source", source).
			Send()
		return
	}

	body, err := reader.GetBody(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error in reading request body"})
		h.logger.Error().
			Err(err).
			Str("source", source).
			Msg("failed to read request body")
		return
	}

	var request models.WithdrawRequest
	if err = json.Unmarshal(body, &request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		h.logger.Error().
			Err(err).
			Str("source", source).
			Msg("failed to unmarshal request")
		return
	}

	request.Username = username

	c, cancel := context.WithTimeout(ctx, h.config.CtxTimeOut)
	defer cancel()

	if err = h.service.Withdraw(c, request); err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidOrderNum):
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": models.ErrInvalidOrderNum.Error()})
			h.logger.Error().
				Err(err).
				Str("source", source).
				Msg("failed to validate order num")
		case errors.Is(err, models.ErrInsufficientFunds):
			ctx.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{"error": models.ErrInsufficientFunds.Error()})
			h.logger.Error().
				Err(err).
				Str("source", source).
				Msg("failed to withdraw")
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error, something went wrong"})
			h.logger.Error().
				Err(err).
				Str("source", source).
				Msg("failed to withdraw")
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "request processed successfully"})
}
