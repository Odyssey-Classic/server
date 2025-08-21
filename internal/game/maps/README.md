# Maps Package

The `maps` package provides data structures and algorithms for managing game maps in the Odyssey RPG server.

## Overview

This package implements a tile-based map system where:
- Each map consists of a 17x17 grid of tiles
- Tiles contain graphics organized by z-index (layers)
- Maps support directional movement blocking, warps, triggers, and custom attributes
- Full JSON serialization/deserialization support

## Architecture

The package is organized into focused files:
- [`types.go`](./types.go) - Core data types and constants
- [`tile.go`](./tile.go) - Tile struct and methods for graphics management
- [`map.go`](./map.go) - Map struct and methods for map operations
- [`maps_test.go`](./maps_test.go) - Comprehensive test suite

## Core Data Structures

### Map
Represents a complete game map with metadata and a 17x17 tile grid. See [`map.go`](./map.go) for the complete structure and methods.

Key fields:
- `ID` - Unique string identifier
- `Name` - Human-readable name
- `Tags` - Array of strings for searching/organizing
- `Attributes` - Custom key/value pairs for gameplay
- `Tiles` - 17x17 grid of tiles
- `Links` - Connections to adjacent maps

### Tile
Represents a single tile with all its properties including graphics, blocking, warps, and triggers. See [`tile.go`](./tile.go) for the complete structure and methods.

Key fields:
- `Passable` - Boolean for basic movement
- `BlockedDirections` - Granular directional blocking
- `Graphics` - Map of z-index to graphics (enforces uniqueness)
- `Warp` - Optional teleport destination
- `Trigger` - Optional script trigger
- `Attributes` - Custom key/value pairs

### Graphics System
Graphics are stored as a map of z-index to `Graphic` objects:
- **Z-index 0 and below**: Rendered beneath player and dynamic objects
- **Z-index 1 and above**: Rendered above player and dynamic objects
- Each z-index can only contain one graphic per tile
- Supports negative z-indexes for underground/background elements

See [`types.go`](./types.go) for the `Graphic` struct definition.

## Key Features

### Directional Blocking
Tiles can block movement in specific cardinal directions (North=0, East=1, South=2, West=3). See [`types.go`](./types.go) for the `DirectionalBlock` struct.

### Warp System
Tiles can contain warp destinations to teleport players to other maps. See [`types.go`](./types.go) for the `WarpDestination` struct.

### Trigger System
Tiles can contain trigger strings for future scripting system integration.

## Usage Examples

For usage examples and patterns, see the comprehensive test suite in [`maps_test.go`](./maps_test.go), which demonstrates:
- Creating and configuring maps
- Working with tiles and graphics
- Checking passability and directional blocking
- Using warps and triggers
- JSON serialization/deserialization

## API Reference

### Map Methods
See [`map.go`](./map.go) for complete method documentation:
- `NewMap(id, name string) *Map` - Creates a new map with default values
- `GetTile(x, y int) (*Tile, error)` - Retrieves a tile at coordinates
- `SetTile(x, y int, tile Tile) error` - Sets a tile at coordinates
- `IsPassable(x, y int, fromDirection Direction) (bool, error)` - Checks if movement is allowed

### Tile Methods
See [`tile.go`](./tile.go) for complete method documentation:
- `AddGraphic(zIndex int, graphic Graphic)` - Adds a graphic to the tile
- `RemoveGraphic(zIndex int)` - Removes a graphic from the tile
- `GetGraphic(zIndex int) (Graphic, bool)` - Retrieves a graphic and existence flag
- `HasGraphic(zIndex int) bool` - Checks if a graphic exists at z-index

### Direction Constants
See [`types.go`](./types.go) for direction definitions:
- `North = 0`
- `East = 1` 
- `South = 2`
- `West = 3`

## Testing

The package includes comprehensive unit tests using the `stretchr/testify/suite` framework:

```bash
go test ./internal/game/maps -v
```

Tests cover all functionality including tile operations, graphics management, passability, serialization, and error handling. See [`maps_test.go`](./maps_test.go) for the complete test suite.

## Design Decisions

1. **17x17 Grid**: Fixed size for consistent memory usage and predictable performance
2. **Map-based Graphics**: Using `map[int]Graphic` enforces z-index uniqueness and supports negative values
3. **Enum Directions**: Using integer constants (0-3) for efficient comparisons
4. **Embedded Layers**: Graphics layers exist within tiles rather than as separate entities
5. **Version Tracking**: Automatic version incrementing on tile modifications
6. **ISO8601 Timestamps**: Human-readable timestamp serialization for database storage
7. **File Separation**: Each type with methods lives in its own file for better encapsulation
