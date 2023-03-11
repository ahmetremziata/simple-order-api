package response

type Order struct {
	OrderNumber  string  `json:"orderNumber"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	TotalAmount  float32 `json:"totalAmount"`
	Address      string  `json:"string"`
	City         string  `json:"city"`
	District     string  `json:"district"`
	CurrencyCode string  `json:"currencyCode"`
	StatusId     int     `json:"statusId"`
}
