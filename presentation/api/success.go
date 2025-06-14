package api

type Success struct {
	Messages []string `json:"messages"`
}

func NewSuccess(messages []string) *Success {
	return &Success{
		Messages: messages,
	}
}
