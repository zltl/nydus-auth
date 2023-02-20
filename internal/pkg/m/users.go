package m

import "github.com/zltl/nydus-auth/pkg/id"

// UserNewReq is the request body for the /user/new route.
type UserNewReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserNewRes struct {
	Error string `json:"error"`
	Id    id.ID  `json:"id"`
}
