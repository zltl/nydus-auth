package m

type AuthorizeLoginReq struct {
	Username    string `form:"username"`
	Password    string `form:"password"`
	ClientId    string `form:"client_id"`
	RedirectURI string `form:"redirect_uri"`
	Scope       string `form:"scope"`
	State       string `form:"state"`

	CSRFToken string `form:"csrf_token"`
}
