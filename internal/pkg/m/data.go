package m

type StoreData struct {
	Ts   int64  `json:"ts"`
	Data string `json:"data"`
	From string `json:"from"` // the hostname that set this data
}
