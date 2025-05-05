package models

type PreferenceRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type PaymentResponse struct {
	Id              int    `json:"id"`
	Status          string `json:"status"`
	DateApproved    string `json:"date_approved"`
	PaymentMethodId string `json:"payment_method_id"`
	PaymentTypeId   string `json:"payment_type_id"`
}

type CreatePreferenceResponse struct {
	Id        string `json:"id"`
	InitPoint string `json:"init_point"`
}