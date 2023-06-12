package inputs

import (
	"github.com/google/uuid"
)

type SendMessage struct {
	ReceiverID *uuid.UUID
	Content    string
}
