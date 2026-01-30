package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// AuthRateLimiter limits auth endpoints to 5 requests per 15 minutes
func AuthRateLimiter() echo.MiddlewareFunc {
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      5,                   // 5 requests
				Burst:     5,                   // burst of 5
				ExpiresIn: 15 * time.Minute,    // per 15 minutes
			},
		),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			// Use IP address as identifier
			return c.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(429, map[string]interface{}{
				"error": map[string]string{
					"message": "Too many requests, please try again later",
				},
			})
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(429, map[string]interface{}{
				"error": map[string]string{
					"message": "Too many requests, please try again later",
				},
			})
		},
	}
	return middleware.RateLimiterWithConfig(config)
}

// GeneralRateLimiter limits general endpoints to 100 requests per minute
func GeneralRateLimiter() echo.MiddlewareFunc {
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      100,              // 100 requests
				Burst:     100,              // burst of 100
				ExpiresIn: 1 * time.Minute,  // per minute
			},
		),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			// Use IP address as identifier
			return c.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(429, map[string]interface{}{
				"error": map[string]string{
					"message": "Too many requests, please try again later",
				},
			})
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(429, map[string]interface{}{
				"error": map[string]string{
					"message": "Too many requests, please try again later",
				},
			})
		},
	}
	return middleware.RateLimiterWithConfig(config)
}
