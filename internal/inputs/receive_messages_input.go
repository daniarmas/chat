package inputs

type ReceiveMessagesInput struct {
	ChatId string `json:"chat_id"`
	UserId string `json:"user_id"`
}

func (in *ReceiveMessagesInput) Sanitize() {
}

func (in ReceiveMessagesInput) Validate() error {
	return nil
}
