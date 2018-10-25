package messages

// Message contains message info
type Message struct {
	Text     string
	Username string
	ID       uint64 `json:"id,uint64"`
}
