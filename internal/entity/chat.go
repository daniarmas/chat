package entity

import (
	"encoding/json"
	"time"
)

type Chat struct {
	ID         string    `json:"id" redis:"id"` // Unique identifier for the message
	Name       string    `json:"name" redis:"name"`
	CreateTime time.Time `json:"create_time" redis:"create_time"`
}

func (m *Chat) MarshalBinary() ([]byte, error) {
	// serialize the message to JSON
	serialized, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	// convert the serialized JSON to binary format
	binaryData := make([]byte, len(serialized))
	copy(binaryData, serialized)

	return binaryData, nil
}
