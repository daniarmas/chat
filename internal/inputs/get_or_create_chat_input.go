package inputs

type GetOrCreateChatInput struct {
	ReceiverId string `json:"receiverId"`
}

func (in *GetOrCreateChatInput) Sanitize() {
}

func (in GetOrCreateChatInput) Validate() error {
	return nil
}
