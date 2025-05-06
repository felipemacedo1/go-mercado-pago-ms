package routes

import (
	"go-mercado-pago-ms/config"
	"go-mercado-pago-ms/controllers"
	"go-mercado-pago-ms/services"

	"github.com/labstack/echo/v4"
	"go-mercado-pago-ms/utils"
	"net/http"
)

func SetupRoutes(e *echo.Echo, cfg config.Config) {
	service := services.MercadoPagoService{Config: cfg}
	controller := controllers.PaymentController{Service: service}

	// Mock JWT Authentication Middleware
	mockAuthMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				utils.Logger.Warn("Request without Authorization header")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}
			// In a real application, you would validate the JWT here.
			return next(c)
		}
	}

	e.POST("/checkout", controller.CreateCheckout, mockAuthMiddleware)
	e.POST("/webhook", controller.ProcessWebhook)
	e.GET("/payment/:id", controller.GetPayment, mockAuthMiddleware)
}
