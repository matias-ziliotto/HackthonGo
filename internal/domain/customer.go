package domain

type Customer struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Situation string `json:"situation"`
}

type CustomerTotalByConditionDTO struct {
	Situation string  `json:"situation"`
	Total     float64 `json:"total"`
}

type CustomerCheaperProductDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
