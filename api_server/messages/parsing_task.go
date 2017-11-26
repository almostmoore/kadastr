package messages

type ParsingTask struct {
	Quarter string `json:"quarter"`
	TextStatus string `json:"text_status"`
	Status int64 `json:"status"`
}
