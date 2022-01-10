package domain

type Product struct {
	Id          int     `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type ProductMostSelledDTO struct {
	Description string  `json:"description"`
	Total       float64 `json:"total"`
}
