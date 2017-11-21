package messages

type AddParsingTaskAnswer struct {
	NotAdded []string `json:"not_added"`
	NotSent []string `json:"not_sent"`
	Added []string `json:"added"`
}
