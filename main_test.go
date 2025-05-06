package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-mercado-pago-ms/cmd"
	"go-mercado-pago-ms/config"
	"go-mercado-pago-ms/models"
	"go-mercado-pago-ms/services"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMercadoPagoService struct {
	mock.Mock
}

func (m *MockMercadoPagoService) CreatePreference(req models.PreferenceRequest) (models.CreatePreferenceResponse, error) {
	args := m.Called(req)
	return args.Get(0).(models.CreatePreferenceResponse), args.Error(1)
}

func (m *MockMercadoPagoService) GetPayment(paymentId string) (models.PaymentResponse, error) {
	args := m.Called(paymentId)
	return args.Get(0).(models.PaymentResponse), args.Error(1)
}

func (m *MockMercadoPagoService) ProcessWebhook(notificationData map[string]interface{}) {
	m.Called(notificationData)
}

func TestMain(t *testing.T) {
	os.Setenv("SERVER_PORT", "8081")
	e := cmd.Init()
	assert.NotNil(t, e)
}

func TestCreateCheckout(t *testing.T) {
	os.Setenv("SERVER_PORT", "8081")
	e := cmd.Init()

	mockService := new(MockMercadoPagoService)
	mockService.On("CreatePreference", mock.Anything).Return(models.CreatePreferenceResponse{Id: "123", InitPoint: "http://test.com"}, nil)

	var cfg config.Config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	paymentService := services.MercadoPagoService{Config: cfg}
	
	
	controller := cmd.GetPaymentController(paymentService)
	e.POST("/checkout", controller.CreateCheckout)

	preference := models.PreferenceRequest{
		Title:       "Test Item",
		Description: "Test Description",
		Price:       100.00,
	}
	
	payload, _ := json.Marshal(preference)

	req := httptest.NewRequest(http.MethodPost, "/checkout", bytes.NewBuffer(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controller.CreateCheckout(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp models.CreatePreferenceResponse
		json.Unmarshal(rec.Body.Bytes(), &resp)

		assert.NotEmpty(t, resp.Id)
		assert.NotEmpty(t, resp.InitPoint)
	}
}

func TestConfigLoad(t *testing.T){
	os.Setenv("SERVER_PORT", "8081")
	os.Setenv("MP_ACCESS_TOKEN", "TEST-123")
	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.Equal(t, "8081", cfg.ServerPort)
	assert.Equal(t, "TEST-123", cfg.MPAccessToken)

	os.Unsetenv("MP_ACCESS_TOKEN")
	os.Unsetenv("SERVER_PORT")
	_, err = config.Load()
	assert.Error(t, err)
}