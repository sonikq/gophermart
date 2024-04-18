package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sonikq/gophermart/internal/models"
	"github.com/sonikq/gophermart/pkg/auth"
	"github.com/sonikq/gophermart/pkg/logger"
	"net/http"
	"strings"
)

func IsAuthorized(log *logger.Logger, secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const source = "middleware.IsAuthorized"
		rawToken := c.GetHeader("Authorization")
		if rawToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{
					"error": models.ErrNotAuthenticated.Error(),
				})
			log.Error().
				Str("source", source).
				Str("error", "Missing Authorization header").
				Send()
			return
		}

		s := strings.Split(rawToken, " ")
		if len(s) != 2 || s[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{
					"error": models.ErrNotAuthenticated.Error(),
				})
			log.Error().
				Str("source", source).
				Str("error", "invalid token").
				Send()
			return
		}

		tkn := s[1]

		username, err := auth.GetUsername(tkn, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": models.ErrNotAuthenticated.Error(),
			})
			log.Error().
				Str("source", source).
				Str("error", "invalid token").
				Send()
		}

		c.Set("username", username)
	}
}
