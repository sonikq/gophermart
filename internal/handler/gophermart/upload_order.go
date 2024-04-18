package gophermart

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sonikq/gophermart/internal/models"
	"github.com/sonikq/gophermart/pkg/reader"
	"net/http"
)

func (h *Handler) UploadOrder(ctx *gin.Context) {
	const source = "handler.UploadOrder"

	username, err := getUsername(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": models.ErrNotAuthenticated.Error()})
		h.logger.Error().
			Err(err).
			Str("source", source).
			Send()
		return
	}

	if ctx.GetHeader(contentTypeHeaderKey) != contentTypeTextPlain {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid type of content"})
		h.logger.Error().
			Str("error", "invalid content type").
			Str("source", source).
			Send()
		return
	}

	bodyBytes, err := reader.GetBody(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error in reading request body"})
		h.logger.Error().
			Err(err).
			Str("source", source).
			Msg("failed to read request body")
		return
	}

	c, cancel := context.WithTimeout(ctx, h.config.CtxTimeOut)
	defer cancel()

	if err = h.service.UploadOrder(c, string(bodyBytes), username); err != nil {
		var (
			statusCode int
			userMsg    string
			logMsg     = "failed to upload order"
		)

		switch {
		case errors.Is(err, models.ErrOrderAlreadyUploadedByAnotherUser):
			statusCode = http.StatusConflict
			userMsg = models.ErrOrderAlreadyUploadedByAnotherUser.Error()

		case errors.Is(err, models.ErrOrderAlreadyUploadedByThisUser):
			statusCode = http.StatusOK
			userMsg = models.ErrOrderAlreadyUploadedByThisUser.Error()

		case errors.Is(err, models.ErrInvalidOrderNum):
			statusCode = http.StatusUnprocessableEntity
			userMsg, logMsg = models.ErrInvalidOrderNum.Error(), "failed to validate order num"

		case errors.Is(err, models.ErrNotUniqueOrderNum):
			statusCode = http.StatusBadRequest
			statusCode = http.StatusUnprocessableEntity
			userMsg, logMsg = models.ErrNotUniqueOrderNum.Error(), "the order number cannot be repeated"

		default:
			statusCode = http.StatusInternalServerError
			userMsg = "internal server error, something went wrong"
		}

		ctx.AbortWithStatusJSON(statusCode, gin.H{"error": userMsg})
		h.logger.Error().Err(err).Str("source", source).Msg(logMsg)
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "New order number accepted for processing"})
}
