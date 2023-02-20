package m

// UserNewReq is the request body for the POST /auth/user route
type SignUpReq struct {
	Username  string `form:"username"`
	Password  string `form:"password"`
	CSRFToken string `form:"csrf_token"`
}
