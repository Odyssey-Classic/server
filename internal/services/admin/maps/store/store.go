package store

import (
	gamemaps "github.com/Odyssey-Classic/server/internal/game/maps"
)

// MapStore abstracts persistence for game maps. Implementations should provide
// durable storage semantics appropriate for the environment (e.g., file-based
// JSON, database, etc.).
type MapStore interface {
	// Create creates a new Map with the given name and returns the fully
	// initialized Map (including assigned ID and timestamps).
	Create(name string) (*gamemaps.Map, error)

	// Get retrieves a Map by its ID.
	Get(id int) (*gamemaps.Map, error)

	// Update persists the provided Map (matching by ID).
	Update(m *gamemaps.Map) error

	// Delete removes a Map by its ID.
	Delete(id int) error

	// List returns all maps, optionally filtered by a case-insensitive
	// substring match on name when query is non-empty.
	List(query string) ([]*gamemaps.Map, error)
}
