package controllers

import (
	"go-mercado-pago-ms/models"
	"go-mercado-pago-ms/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PaymentController struct {
	Service services.MercadoPagoService
}

func (pc *PaymentController) CreateCheckout(c echo.Context) error {
	var req models.PreferenceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := pc.Service.CreatePreference(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (pc *PaymentController) ProcessWebhook(c echo.Context) error {
	var body map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
	}

	go pc.Service.ProcessWebhook(body)

	return c.NoContent(http.StatusOK)
}

func (pc *PaymentController) GetPayment(c echo.Context) error {
	paymentID := c.Param("id")
	payment, err := pc.Service.GetPayment(paymentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, payment)
}


