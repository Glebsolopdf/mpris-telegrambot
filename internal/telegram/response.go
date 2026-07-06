package telegram

import "encoding/json"

type response struct {
	OK          bool            `json:"ok"`
	Description string          `json:"description"`
	Parameters  parameters      `json:"parameters"`
	Result      json.RawMessage `json:"result"`
}

type parameters struct {
	RetryAfter int `json:"retry_after"`
}
