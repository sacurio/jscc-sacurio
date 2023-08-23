package message

// Message defines the basic structure of a message.
type Message struct {
	Queue   string `json:"queue"`
	Content string `json:"content"`
}
