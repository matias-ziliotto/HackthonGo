package domain

type Invoice struct {
	Id          int     `json:"id"`
	Customer_id int     `json:"customer_id"`
	Datetime    string  `json:"datetime"`
	Total       float64 `json:"total"`
}

type InvoiceDTO struct {
	Id       int      `json:"id"`
	Customer Customer `json:"customer"`
	Datetime string   `json:"datetime"`
	Total    float64  `json:"total"`
}

type InvoiceTotalDTO struct {
	Id    int     `json:"id"`
	Total float64 `json:"total"`
}
