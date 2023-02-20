package m

// ErrRes is the response body for an error.
type ErrRes struct {
	Error string `json:"error"`
	Msg   string `json:"msg"`
}
