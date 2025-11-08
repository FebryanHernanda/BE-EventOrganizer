package middleware

import (
	"strings"

	"github.com/FebryanHernanda/BE-EventOrganizer/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func JWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			logrus.Warn("Missing Authorization header")
			response.Error(ctx, "invalid token format, expected `Bearer <token>'", 401, nil)
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			logrus.Warn("Invalid token format (missing `Bearer `)")
			response.Error(ctx, "invalid token format, expected `Bearer <token>`", 401, nil)
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			logrus.WithError(err).Warn("Invalid token")
			response.Error(ctx, "invalid or expired token", 401, err.Error())
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx.Set("user_id", claims["user_id"])
			ctx.Set("role", claims["role"])
			ctx.Next()
			return
		}

		logrus.Warn("Invalid token claims")
		response.Error(ctx, "invalid token claims", 401, nil)
		ctx.Abort()

	}
}
