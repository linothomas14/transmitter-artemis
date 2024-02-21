package entity

type Queue struct {
	MessageID string `json:"message_id"`
	To        string `json:"to"`
	Type      string `json:"type"`
	Text      struct {
		PreviewURL bool   `json:"preview_url"`
		Body       string `json:"body"`
	}
}
