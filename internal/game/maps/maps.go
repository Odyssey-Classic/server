// Package maps provides data structures and algorithms for managing game maps in the Odyssey RPG server.
//
// This package implements a tile-based map system where:
// - Each map consists of a 17x17 grid of tiles
// - Tiles contain graphics organized by z-index (layers)
// - Maps support directional movement blocking, warps, triggers, and custom attributes
// - Full JSON serialization/deserialization support
//
// The package is organized into several files:
// - types.go: Core data types and constants
// - tile.go: Tile struct and methods
// - map.go: Map struct and methods
// - maps_test.go: Comprehensive test suite
package maps
