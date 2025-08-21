package maps

// Direction represents a cardinal direction using enum values.
type Direction int

const (
	North Direction = 0
	East  Direction = 1
	South Direction = 2
	West  Direction = 3
)

// Delta returns the X, Y coordinate changes for moving in this direction.
// North: (0, -1), East: (1, 0), South: (0, 1), West: (-1, 0)
func (d Direction) Delta() (int, int) {
	switch d {
	case North:
		return 0, -1
	case East:
		return 1, 0
	case South:
		return 0, 1
	case West:
		return -1, 0
	default:
		return 0, 0 // Invalid direction
	}
}

// DirectionalBlock represents blocking rules for a direction.
type DirectionalBlock struct {
	Direction     Direction `json:"direction"`
	BlockInbound  bool      `json:"block_inbound"`
	BlockOutbound bool      `json:"block_outbound"`
}

// Graphic represents a graphic on a tile with extensible properties.
type Graphic struct {
	GraphicID  int               `json:"graphic_id"`
	Properties map[string]string `json:"properties,omitempty"`
}

// WarpDestination represents a warp target location.
type WarpDestination struct {
	MapID int `json:"map_id"`
	X     int `json:"x"`
	Y     int `json:"y"`
}

// MapLinks represents links to other maps by direction.
type MapLinks struct {
	North int `json:"north,omitempty"`
	East  int `json:"east,omitempty"`
	South int `json:"south,omitempty"`
	West  int `json:"west,omitempty"`
}
