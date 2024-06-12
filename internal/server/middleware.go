package server

import (
	"app/internal/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

// Checks if the token is valid and extracts user id from it
func (server *Server) authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return

		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		token := fields[1]
		if err := util.VerifyToken(token, []byte(server.Config.SecretKey)); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		tokenClaims, ok := util.ExtractClaims(token, []byte(server.Config.SecretKey))
		if !ok {
			err := errors.New("invalid token")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		userID, ok := tokenClaims["user_id"]
		if !ok {
			err := errors.New("invalid token")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
