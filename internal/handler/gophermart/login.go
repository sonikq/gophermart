package gophermart

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sonikq/gophermart/internal/models"
	"github.com/sonikq/gophermart/pkg/auth"
	"github.com/sonikq/gophermart/pkg/reader"
	"net/http"
)

func (h *Handler) Login(ctx *gin.Context) {
	const source = "handler.Login"

	if ctx.GetHeader(contentTypeHeaderKey) != contentTypeJSON {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid type of content"})
		h.logger.Error().
			Err(errors.New("invalid content type")).
			Str("source", source).
			Send()
	}

	bodyBytes, err := reader.GetBody(ctx.Request.Body)
	if err != nil {
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error in reading request body"})
			h.logger.Error().
				Err(err).
				Str("source", source).
				Msg("failed to read request body")
			return
		}
	}

	var request models.LoginRequest
	if err = json.Unmarshal(bodyBytes, &request); err != nil {
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			h.logger.Error().
				Err(err).
				Str("source", source).
				Msg("failed to parse request body")
			return
		}
	}

	if err = request.Validate(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.logger.Error().
			Err(err).
			Str("source", source).
			Msg("failed to validate request")
		return
	}

	c, cancel := context.WithTimeout(ctx, h.config.CtxTimeOut)
	defer cancel()

	if err = h.service.Login(c, request); err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Invalid credentials"})
			h.logger.Error().
				Err(err).
				Str("source", source).
				Send()
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal error, something went wrong"})
		h.logger.Error().
			Err(err).
			Str("source", source).
			Send()
		return
	}

	// Authentication
	token, err := auth.GenerateToken(h.config.TokenExpiration, h.config.TokenSecretKey, request.Login)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal error, something went wrong"})
		h.logger.Error().
			Err(err).
			Str("source", source).
			Send()
		return
	}

	ctx.Header("Authorization", "Bearer "+token)
	ctx.JSON(http.StatusOK, gin.H{"message": "The user is successfully authenticated"})
}
