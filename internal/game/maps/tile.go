package maps

// Tile represents a single tile with all its properties.
// Graphics are stored as a map of z-index to Graphic, enforcing uniqueness
// and allowing negative/positive z-indexes.
// Z-index 0 and below: rendered beneath player and dynamic objects
// Z-index 1 and above: rendered above player and dynamic objects
type Tile struct {
	Passable          bool               `json:"passable"`
	BlockedDirections []DirectionalBlock `json:"blocked_directions,omitempty"`
	Graphics          map[int]Graphic    `json:"graphics,omitempty"`
	Warp              *WarpDestination   `json:"warp,omitempty"`
	Trigger           string             `json:"trigger,omitempty"`
	Attributes        map[string]string  `json:"attributes,omitempty"`
}

// AddGraphic adds a graphic at the specified z-index to this tile.
func (t *Tile) AddGraphic(zIndex int, graphic Graphic) {
	if t.Graphics == nil {
		t.Graphics = make(map[int]Graphic)
	}
	t.Graphics[zIndex] = graphic
}

// RemoveGraphic removes a graphic at the specified z-index from this tile.
func (t *Tile) RemoveGraphic(zIndex int) {
	if t.Graphics != nil {
		delete(t.Graphics, zIndex)
	}
}

// GetGraphic retrieves a graphic at the specified z-index from this tile.
func (t *Tile) GetGraphic(zIndex int) (Graphic, bool) {
	if t.Graphics == nil {
		return Graphic{}, false
	}
	graphic, exists := t.Graphics[zIndex]
	return graphic, exists
}

// HasGraphic checks if a graphic exists at the specified z-index.
func (t *Tile) HasGraphic(zIndex int) bool {
	_, exists := t.GetGraphic(zIndex)
	return exists
}
