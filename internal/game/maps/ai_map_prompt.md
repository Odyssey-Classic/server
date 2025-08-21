# AI Agent Prompt: Map Data Structures & Algorithms for Odyssey RPG

## Context
You are an AI agent tasked with designing and implementing map data structures and algorithms for Odyssey, a tile-based RPG. The game uses a 17x17 tile map with multiple layers (e.g., ground, objects, effects).

## Requirements

- The map should support multiple layers (minimum: ground, objects, effects).
- Each layer consists of a 17x17 grid of tiles.
- Tiles may contain properties such as type, passability, events, and a reference to which graphics to use for each tile on each layer.
- Efficient access and modification of tiles is required.
- Support for serialization/deserialization (saving/loading maps).
- Algorithms for pathfinding, tile updates, and layer management.
- Each map can have links to other maps for each cardinal direction (north, south, east, west).

## Tasks

1. Design data structures to represent the map, layers, and tiles.
2. Implement algorithms for:
    - Accessing and modifying tiles.
    - Managing multiple layers.
    - Pathfinding (e.g., A* or Dijkstra’s algorithm).
    - Serializing and deserializing map data.
3. Provide code samples and explanations for each component.

## Output Format

- Use clear, well-documented code.
- Include comments explaining design decisions.
- Provide example usage for each major function or class.
- Build unit tests using the stretchr/testify/suite package.

## Constraints

- Focus on clarity, efficiency, and extensibility.

## Example Tile Properties

- `passable`: boolean
- `event`: optional (e.g., triggers, NPCs)
- `graphic`: string or object reference (e.g., sprite ID, tileset index, or asset path)
