package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/config"
)

type JWTClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func JWTAuth(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]string{
						"message": "Missing authorization header",
					},
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]string{
						"message": "Invalid authorization header format",
					},
				})
			}

			tokenString := parts[1]

			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]string{
						"message": "Invalid or expired token",
					},
				})
			}

			claims, ok := token.Claims.(*JWTClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]string{
						"message": "Invalid token claims",
					},
				})
			}

			// Set user info in context
			c.Set("userId", claims.UserID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}

func GetUserID(c echo.Context) string {
	userID, _ := c.Get("userId").(string)
	return userID
}
