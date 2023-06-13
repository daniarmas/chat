package inputs

import (
	"github.com/google/uuid"
)

type SendMessage struct {
	ChatID  *uuid.UUID
	Content string
}
