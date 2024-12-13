package platform

type GoogleUser struct {
	ID      string `json:"sub"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
