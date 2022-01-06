package domain

type Sale struct {
	Id         int     `json:"id"`
	Invoice_id int     `json:"invoice_id"`
	Product_id int     `json:"product_id"`
	Quantity   float64 `json:"quantity"`
}

type SaleDTO struct {
	Id       int     `json:"id"`
	Invoice  Invoice `json:"invoice"`
	Product  Product `json:"product"`
	Quantity float64 `json:"quantity"`
}
