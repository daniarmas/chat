package inputs

import (
	"github.com/google/uuid"
)

type GetOrCreateChatInput struct {
	OtherUserId *uuid.UUID `json:"otherUserId"`
}

func (in *GetOrCreateChatInput) Sanitize() {
}

func (in GetOrCreateChatInput) Validate() error {
	return nil
}
