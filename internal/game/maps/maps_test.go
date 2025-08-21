package maps

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type MapSuite struct {
	suite.Suite
	m *Map
}

func (s *MapSuite) SetupTest() {
	s.m = NewMap(1, "TestMap")
	s.m.Tags = []string{"test", "rpg"}
	s.m.Attributes = map[string]string{"difficulty": "easy"}
	s.m.Links = MapLinks{North: 2}
}

func (s *MapSuite) TestNewMap() {
	m := NewMap(100, "NewTestMap")
	s.Equal(100, m.ID)
	s.Equal("NewTestMap", m.Name)
	s.Equal(1, m.Version)
	s.NotZero(m.LastUpdated)
	s.Empty(m.Tags)
	s.NotNil(m.Attributes)
}

func (s *MapSuite) TestTileAccessAndModification() {
	// Test basic tile access
	tile, err := s.m.GetTile(5, 5)
	s.Require().NoError(err)
	s.False(tile.Passable) // Default should be false

	// Test tile modification
	newTile := Tile{
		Passable: true,
		Graphics: map[int]Graphic{
			0:  {GraphicID: 100, Properties: map[string]string{"variant": "dark"}},
			1:  {GraphicID: 200},
			-1: {GraphicID: 300},
		},
		Attributes: map[string]string{"biome": "forest"},
	}

	oldVersion := s.m.Version
	err = s.m.SetTile(5, 5, newTile)
	s.Require().NoError(err)
	s.Equal(oldVersion+1, s.m.Version) // Version should increment

	// Verify the tile was set correctly
	retrievedTile, err := s.m.GetTile(5, 5)
	s.Require().NoError(err)
	s.True(retrievedTile.Passable)
	s.Equal(100, retrievedTile.Graphics[0].GraphicID)
	s.Equal("dark", retrievedTile.Graphics[0].Properties["variant"])
	s.Equal(200, retrievedTile.Graphics[1].GraphicID)
	s.Equal(300, retrievedTile.Graphics[-1].GraphicID)
	s.Equal("forest", retrievedTile.Attributes["biome"])
}

func (s *MapSuite) TestTileOutOfBounds() {
	_, err := s.m.GetTile(-1, 5)
	s.Error(err)

	_, err = s.m.GetTile(17, 5)
	s.Error(err)

	_, err = s.m.GetTile(5, -1)
	s.Error(err)

	_, err = s.m.GetTile(5, 17)
	s.Error(err)
}

func (s *MapSuite) TestGraphicManagement() {
	graphic := Graphic{
		GraphicID:  400,
		Properties: map[string]string{"hardness": "high"},
	}

	// Get the tile and add graphic to it
	tile, err := s.m.GetTile(3, 3)
	s.Require().NoError(err)

	tile.AddGraphic(2, graphic)
	err = s.m.SetTile(3, 3, *tile)
	s.Require().NoError(err)

	// Verify graphic was added
	tile, err = s.m.GetTile(3, 3)
	s.Require().NoError(err)
	s.Equal(400, tile.Graphics[2].GraphicID)
	s.Equal("high", tile.Graphics[2].Properties["hardness"])
	s.True(tile.HasGraphic(2))

	retrievedGraphic, exists := tile.GetGraphic(2)
	s.True(exists)
	s.Equal(400, retrievedGraphic.GraphicID)

	// Remove graphic
	tile.RemoveGraphic(2)
	err = s.m.SetTile(3, 3, *tile)
	s.Require().NoError(err)

	tile, err = s.m.GetTile(3, 3)
	s.Require().NoError(err)
	s.False(tile.HasGraphic(2))
	_, exists = tile.GetGraphic(2)
	s.False(exists)
}

func (s *MapSuite) TestPassability() {
	// Set up a passable tile
	passableTile := Tile{Passable: true}
	err := s.m.SetTile(0, 0, passableTile)
	s.Require().NoError(err)

	passable, err := s.m.IsPassable(0, 0, North)
	s.Require().NoError(err)
	s.True(passable)

	// Set up a tile with directional blocking
	blockedTile := Tile{
		Passable: true,
		BlockedDirections: []DirectionalBlock{
			{Direction: North, BlockInbound: true, BlockOutbound: false},
		},
	}
	err = s.m.SetTile(1, 1, blockedTile)
	s.Require().NoError(err)

	passable, err = s.m.IsPassable(1, 1, North)
	s.Require().NoError(err)
	s.False(passable)

	passable, err = s.m.IsPassable(1, 1, East)
	s.Require().NoError(err)
	s.True(passable)
}

func (s *MapSuite) TestWarpFunctionality() {
	warpTile := Tile{
		Passable: true,
		Warp: &WarpDestination{
			MapID: 10,
			X:     5,
			Y:     10,
		},
	}

	err := s.m.SetTile(8, 8, warpTile)
	s.Require().NoError(err)

	tile, err := s.m.GetTile(8, 8)
	s.Require().NoError(err)
	s.NotNil(tile.Warp)
	s.Equal(10, tile.Warp.MapID)
	s.Equal(5, tile.Warp.X)
	s.Equal(10, tile.Warp.Y)
}

func (s *MapSuite) TestTriggerFunctionality() {
	triggerTile := Tile{
		Passable: true,
		Trigger:  "open_chest_script",
	}

	err := s.m.SetTile(2, 2, triggerTile)
	s.Require().NoError(err)

	tile, err := s.m.GetTile(2, 2)
	s.Require().NoError(err)
	s.Equal("open_chest_script", tile.Trigger)
}

func (s *MapSuite) TestSerialization() {
	// Set up a complex tile for testing
	complexTile := Tile{
		Passable: true,
		BlockedDirections: []DirectionalBlock{
			{Direction: North, BlockInbound: true, BlockOutbound: false},
		},
		Graphics: map[int]Graphic{
			0:  {GraphicID: 500},
			1:  {GraphicID: 600, Properties: map[string]string{"color": "red"}},
			-1: {GraphicID: 700},
		},
		Warp:       &WarpDestination{MapID: 20, X: 0, Y: 0},
		Trigger:    "entrance_script",
		Attributes: map[string]string{"danger_level": "low"},
	}

	err := s.m.SetTile(10, 10, complexTile)
	s.Require().NoError(err)

	// Test JSON serialization
	data, err := json.Marshal(s.m)
	s.Require().NoError(err)

	// Test JSON deserialization
	var m2 Map
	err = json.Unmarshal(data, &m2)
	s.Require().NoError(err)

	// Verify the deserialized map matches the original
	s.Equal(s.m.Name, m2.Name)
	s.Equal(s.m.Version, m2.Version)
	s.Equal(s.m.Tags, m2.Tags)
	s.Equal(s.m.Attributes, m2.Attributes)
	s.Equal(s.m.Links, m2.Links)

	// Verify the complex tile was preserved
	tile := m2.Tiles[10][10]
	s.True(tile.Passable)
	s.Len(tile.BlockedDirections, 1)
	s.Equal(North, tile.BlockedDirections[0].Direction)
	s.True(tile.BlockedDirections[0].BlockInbound)
	s.False(tile.BlockedDirections[0].BlockOutbound)

	s.Equal(500, tile.Graphics[0].GraphicID)
	s.Equal(600, tile.Graphics[1].GraphicID)
	s.Equal("red", tile.Graphics[1].Properties["color"])
	s.Equal(700, tile.Graphics[-1].GraphicID)

	s.Equal(20, tile.Warp.MapID)
	s.Equal("entrance_script", tile.Trigger)
	s.Equal("low", tile.Attributes["danger_level"])

	// Verify LastUpdated was serialized/deserialized correctly
	s.WithinDuration(s.m.LastUpdated, m2.LastUpdated, time.Second)
}

func (s *MapSuite) TestDirectionConstants() {
	s.Equal(Direction(0), North)
	s.Equal(Direction(1), East)
	s.Equal(Direction(2), South)
	s.Equal(Direction(3), West)
}

func (s *MapSuite) TestDirectionDelta() {
	// Test North
	dx, dy := North.Delta()
	s.Equal(0, dx)
	s.Equal(-1, dy)

	// Test East
	dx, dy = East.Delta()
	s.Equal(1, dx)
	s.Equal(0, dy)

	// Test South
	dx, dy = South.Delta()
	s.Equal(0, dx)
	s.Equal(1, dy)

	// Test West
	dx, dy = West.Delta()
	s.Equal(-1, dx)
	s.Equal(0, dy)

	// Test invalid direction
	invalidDir := Direction(99)
	dx, dy = invalidDir.Delta()
	s.Equal(0, dx)
	s.Equal(0, dy)
}

func TestMapSuite(t *testing.T) {
	suite.Run(t, new(MapSuite))
}
