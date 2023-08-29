package model

type Chat struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	Msg       string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type ChatError struct {
	Msg      string `json:"error"`
	Status   int    `json:"status"`
	CausedBy string `json:"causedby"`
}

type ContactList struct {
	Username     string `json:"username"`
	LastActivity int64  `json:"last_activity"`
}
