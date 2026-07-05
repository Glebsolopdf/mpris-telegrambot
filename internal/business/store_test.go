package business

import (
	"path/filepath"
	"testing"
)

func TestStoreRoundTripForUser(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "business_connection.json"))
	state := State{ConnectionID: "bc_1", UserID: 42}

	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}

	loaded, ok, err := store.LoadForUser(42)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || loaded.ConnectionID != state.ConnectionID {
		t.Fatalf("loaded = %+v, ok = %v", loaded, ok)
	}

	if _, ok, err := store.LoadForUser(7); err != nil || ok {
		t.Fatalf("wrong user ok = %v, err = %v", ok, err)
	}

	if err := store.Delete(); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := store.LoadForUser(42); err != nil || ok {
		t.Fatalf("deleted store ok = %v, err = %v", ok, err)
	}
}
