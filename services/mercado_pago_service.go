package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-mercado-pago-ms/config"
	"go-mercado-pago-ms/models"
	"go-mercado-pago-ms/utils"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type MercadoPagoService struct {
	Config config.Config
}

type preferenceItem struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

type preferenceRequest struct {
	Items               []preferenceItem  `json:"items"`
	BackUrls            map[string]string `json:"back_urls"`
	AutoReturn          string            `json:"auto_return"`
	StatementDescriptor string            `json:"statement_descriptor,omitempty"`
}

func (s *MercadoPagoService) CreatePreference(req models.PreferenceRequest) (models.CreatePreferenceResponse, error) {

	preference := preferenceRequest{
		Items: []preferenceItem{
			{
				Title:       req.Title,
				Description: req.Description,
				Quantity:    1,
				UnitPrice:   int(req.Price),
			},
		},
		BackUrls: map[string]string{
			"success": "http://localhost:8080/success",
			"failure": "http://localhost:8080/failure",
			"pending": "http://localhost:8080/pending",
		},
		AutoReturn:          "approved",
		StatementDescriptor: "MERCADOPAGO",
	}

	payload, err := json.Marshal(preference)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error creating preference: %v", err))
		return models.CreatePreferenceResponse{}, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	url := "https://api.mercadopago.com/checkout/preferences"
	reqHttp, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error creating preference: %v", err))
		return models.CreatePreferenceResponse{}, err
	}

	reqHttp.Header.Set("Content-Type", "application/json")
	reqHttp.Header.Set("Authorization", "Bearer "+s.Config.MPAccessToken)

	resp, err := client.Do(reqHttp)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error creating preference: %v", err))
		return models.CreatePreferenceResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		utils.Logger.Error(fmt.Sprintf("Error creating preference: %s", bodyString))
		return models.CreatePreferenceResponse{}, errors.New(fmt.Sprintf("Error creating preference: %s", bodyString))
	}

	var createPreferenceResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&createPreferenceResponse)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error creating preference: %v", err))
		return models.CreatePreferenceResponse{}, err
	}

	response := models.CreatePreferenceResponse{
		Id:        createPreferenceResponse["id"].(string),
		InitPoint: createPreferenceResponse["init_point"].(string),
	}

	return response, nil
}

func (s *MercadoPagoService) GetPayment(paymentId string) (models.PaymentResponse, error) {
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s", paymentId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error getting payment: %v", err))
		return models.PaymentResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+s.Config.MPAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error getting payment: %v", err))
		return models.PaymentResponse{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error reading body payment: %v", err))
		return models.PaymentResponse{}, err
	}

	bodyString := string(bodyBytes)
	if resp.StatusCode != http.StatusOK {
		utils.Logger.Error(fmt.Sprintf("Error getting payment: %v", bodyString))
		return models.PaymentResponse{}, errors.New(fmt.Sprintf("Error getting payment: %v", bodyString))
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(bodyString), &result)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error unmarshalling payment: %v", err))
		return models.PaymentResponse{}, err
	}

	dateApproved, ok := result["date_approved"].(string)
	if !ok {
		utils.Logger.Error(fmt.Sprintf("Error getting date_approved: %v", err))
		return models.PaymentResponse{}, err
	}
	paymentMethodId, ok := result["payment_method_id"].(string)
	if !ok {
		utils.Logger.Error(fmt.Sprintf("Error getting payment_method_id: %v", err))
		return models.PaymentResponse{}, err
	}
	paymentTypeId, ok := result["payment_type_id"].(string)
	if !ok {
		utils.Logger.Error(fmt.Sprintf("Error getting payment_type_id: %v", err))
		return models.PaymentResponse{}, err
	}

	paymentResponse := models.PaymentResponse{
		Id:              int(result["id"].(float64)),
		Status:          result["status"].(string),
		DateApproved:    dateApproved,
		PaymentMethodId: paymentMethodId,
		PaymentTypeId:   paymentTypeId,
	}

	return paymentResponse, nil
}

func (s *MercadoPagoService) ProcessWebhook(notificationData map[string]interface{}) {
	utils.Logger.Info(fmt.Sprintf("Webhook notification received: %v", notificationData))

	// Process the notification data here
	// You can log the data, store it in a database, or perform other actions
	// depending on your needs.

	if resource, ok := notificationData["resource"].(string); ok {
		utils.Logger.Info("resource:", resource)

		parts := strings.Split(resource, "/")
		if len(parts) > 1 {
			paymentId := parts[len(parts)-1]
			utils.Logger.Info("paymentId:", paymentId)
			payment, err := s.GetPayment(paymentId)
			if err != nil {
				utils.Logger.Error(fmt.Sprintf("Error getting payment: %v", err))
			} else {
				utils.Logger.Info(fmt.Sprintf("Payment Info: %v", payment))
			}
		}

	}

	if topic, ok := notificationData["topic"].(string); ok {
		utils.Logger.Info("topic:", topic)
		if topic == "merchant_order" {
			utils.Logger.Info("Processing merchant order notification")
		}
	}

	file, err := os.OpenFile("webhook.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error creating webhook log: %v", err))
		return
	}
	defer file.Close()
	data, err := json.Marshal(notificationData)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error unmarshalling log data: %v", err))
		return
	}
	if _, err := file.WriteString(fmt.Sprintf("%s\n", data)); err != nil {
		utils.Logger.Error(fmt.Sprintf("Error writing log data: %v", err))
		return
	}

}
