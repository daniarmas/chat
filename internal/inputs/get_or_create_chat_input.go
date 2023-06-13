package inputs

import (
	"github.com/google/uuid"
)

type GetOrCreateChatInput struct {
	FirstUserId  *uuid.UUID `json:"firstUserId"`
	SecondUserId *uuid.UUID `json:"secondUserId"`
}

func (in *GetOrCreateChatInput) Sanitize() {
}

func (in GetOrCreateChatInput) Validate() error {
	return nil
}
