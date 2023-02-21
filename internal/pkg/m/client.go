package m

type AuthClient struct {
	Id          string   `json:"client_id"`
	Name        string   `json:"name"`
	Secret      string   `json:"client_secret"`
	Scopes      []string `json:"scope"`
	RedirectUri []string `json:"redirect_uri"`
}
