package item

type Item struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Plaintext   string   `json:"plaintext"`
	Tags        []string `json:"tags"`
	BasePrice   int      `json:"basePrice"`
	TotalPrice  int      `json:"totalPrice"`
	SellPrice   int      `json:"sellPrice"`
	Purchasable bool     `json:"purchasable"`
	ImageURL    string   `json:"imageUrl"`
}
