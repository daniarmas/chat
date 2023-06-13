package inputs

import (
	"github.com/google/uuid"
)

type GetOrCreateChatInput struct {
	ReceiverId *uuid.UUID `json:"receiverId"`
}

func (in *GetOrCreateChatInput) Sanitize() {
}

func (in GetOrCreateChatInput) Validate() error {
	return nil
}
