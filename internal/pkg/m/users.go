package m

import "github.com/zltl/nydus-auth/pkg/id"

type UserPassword struct {
	ID       id.ID  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserNewRes struct {
	Error string `json:"error"`
	Id    id.ID  `json:"id"`
}
