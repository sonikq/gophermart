package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sonikq/gophermart/pkg/logger"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

type gzipReader struct {
	io.ReadCloser
	reader *gzip.Reader
}

func (g *gzipReader) Read(data []byte) (int, error) {
	return g.reader.Read(data)
}

func (g *gzipReader) Close() error {
	return g.reader.Close()
}

func CompressResponse(log *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		source := "middleware.CompressResponse"
		acceptEncoding := ctx.GetHeader("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			ctx.Next()
			return
		}

		gz := gzip.NewWriter(ctx.Writer)
		defer func() {
			if err := gz.Close(); err != nil {
				log.Info().
					Str("source", source).
					Str("error", "error in closing gzip writer").
					Send()
				return
			}
		}()

		ctx.Header("Content-Encoding", "gzip")
		ctx.Writer = &gzipWriter{
			ResponseWriter: ctx.Writer,
			writer:         gz,
		}
		ctx.Next()
	}
}

func DecompressRequest(log *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		source := "middleware.DecompressRequest"
		if !strings.Contains(ctx.GetHeader("Content-Encoding"), "gzip") {
			ctx.Next()
			return
		}
		gz, err := gzip.NewReader(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		defer func() {
			if err = gz.Close(); err != nil {
				log.Info().
					Str("source", source).
					Str("error", "error in closing gzip reader").
					Send()
				return
			}
		}()

		ctx.Request.Body = &gzipReader{
			ReadCloser: ctx.Request.Body,
			reader:     gz,
		}
		ctx.Next()
	}
}
