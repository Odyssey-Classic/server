package maps

import (
	"encoding/json"
	"fmt"
	"time"
)

// Map represents a game map with a 17x17 grid of tiles.
type Map struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Tags        []string          `json:"tags"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	LastUpdated time.Time         `json:"last_updated"`
	Version     int               `json:"version"`
	Tiles       [17][17]Tile      `json:"tiles"`
	Links       MapLinks          `json:"links"`
}

// MarshalJSON customizes serialization of Map to format LastUpdated as ISO8601.
func (m Map) MarshalJSON() ([]byte, error) {
	type Alias Map
	return json.Marshal(&struct {
		LastUpdated string `json:"last_updated"`
		*Alias
	}{
		LastUpdated: m.LastUpdated.Format(time.RFC3339),
		Alias:       (*Alias)(&m),
	})
}

// UnmarshalJSON customizes deserialization of Map to parse LastUpdated from ISO8601.
func (m *Map) UnmarshalJSON(data []byte) error {
	type Alias Map
	aux := &struct {
		LastUpdated string `json:"last_updated"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	parsed, err := time.Parse(time.RFC3339, aux.LastUpdated)
	if err != nil {
		return err
	}
	m.LastUpdated = parsed
	return nil
}

// GetTile returns the tile at the given x, y coordinates.
func (m *Map) GetTile(x, y int) (*Tile, error) {
	if x < 0 || x >= 17 || y < 0 || y >= 17 {
		return nil, fmt.Errorf("coordinates out of range: (%d, %d)", x, y)
	}
	return &m.Tiles[x][y], nil
}

// SetTile updates the tile at the given x, y coordinates.
func (m *Map) SetTile(x, y int, tile Tile) error {
	if x < 0 || x >= 17 || y < 0 || y >= 17 {
		return fmt.Errorf("coordinates out of range: (%d, %d)", x, y)
	}
	m.Tiles[x][y] = tile
	m.LastUpdated = time.Now()
	m.Version++
	return nil
}

// IsPassable checks if a tile is passable and not blocked in a specific direction.
func (m *Map) IsPassable(x, y int, fromDirection Direction) (bool, error) {
	tile, err := m.GetTile(x, y)
	if err != nil {
		return false, err
	}

	if !tile.Passable {
		return false, nil
	}

	// Check directional blocks
	for _, block := range tile.BlockedDirections {
		if block.Direction == fromDirection && (block.BlockInbound || block.BlockOutbound) {
			return false, nil
		}
	}

	return true, nil
}

// NewMap creates a new map with default values.
func NewMap(id int, name string) *Map {
	return &Map{
		ID:          id,
		Name:        name,
		Tags:        make([]string, 0),
		Attributes:  make(map[string]string),
		LastUpdated: time.Now(),
		Version:     1,
		Tiles:       [17][17]Tile{},
		Links:       MapLinks{},
	}
}
