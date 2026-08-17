package champion

type Champion struct {
	ID       string   `json:"id"`
	Key      string   `json:"key"`
	Name     string   `json:"name"`
	Title    string   `json:"title"`
	Blurb    string   `json:"blurb"`
	Tags     []string `json:"tags"`
	ImageURL string   `json:"imageUrl"`
}
