package models

type CEPRequest struct {
	CEP string `json:"cep" binding:"required"`
}

type CEPResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
