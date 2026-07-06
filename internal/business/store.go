package business

import (
	"encoding/json"
	"errors"
	"os"
)

type State struct {
	ConnectionID string `json:"business_connection_id"`
	UserID       int64  `json:"user_id"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
}

type Store struct {
	path string
}

func NewStore(path string) Store {
	return Store{path: path}
}

func (s Store) LoadForUser(userID int64) (State, bool, error) {
	payload, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return State{}, false, nil
	}
	if err != nil {
		return State{}, false, err
	}

	var state State
	if err := json.Unmarshal(payload, &state); err != nil {
		return State{}, false, err
	}
	if state.UserID != userID || state.ConnectionID == "" {
		return State{}, false, nil
	}
	return state, true, nil
}

func (s Store) Save(state State) error {
	payload, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, append(payload, '\n'), 0o600)
}

func (s Store) Delete() error {
	if err := os.Remove(s.path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else {
		return err
	}
}
