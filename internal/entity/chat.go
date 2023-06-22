package entity

import (
	"encoding/json"
	"time"
)

type Chat struct {
	ID           string    `json:"id" redis:"id"` // Unique identifier for the message
	FirstUserId  string    `json:"first_user_id" redis:"first_user_id"`
	SecondUserId string    `json:"second_user_id" redis:"second_user_id"`
	CreateTime   time.Time `json:"create_time" redis:"create_time"`
	UpdateTime   time.Time `json:"update_time" redis:"update_time"`
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
