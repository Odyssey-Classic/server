# AI Agent Prompt: Map Data Structures & Algorithms for Odyssey RPG

## Context
You are an AI agent tasked with designing and implementing map data structures and algorithms for Odyssey, a tile-based RPG. The game uses a 17x17 tile map with multiple layers (e.g., ground, objects, effects).

## Requirements

- Each map contains the following properties:
    - Name: string, the name of the map.
    - Tags: array of strings, used for searching and organizing maps only.
    - Attributes: key/value pairs (map of strings), intended to affect gameplay on the map. No predefined keys; all keys are user-defined.
    - Administrative fields:
        - Last update timestamp: stored as simple data (e.g., integer or time object), but serialized as a human-readable string (e.g., ISO8601).
        - Version number: integer.
- Each layer consists of a 17x17 grid of tiles.
- Each tile contains the following properties:
    - Passability: a boolean for unpassable tiles.
    - Directional blocking: ability to block movement into or out of the tile from any cardinal direction (north, east, south, west only; diagonal directions are not supported). Each block must specify the direction and whether it blocks inbound, outbound, or both.
    - Graphic layers: an array, each with a unique integer z-index, a graphic ID, and an extensible list of properties for future expansion. Only one graphic can exist at each z-index for a tile.
    - Attributes: key/value pairs (map of strings), intended to affect gameplay for the tile. No predefined keys; all keys are user-defined.
- Layers with z-index 0 and below are rendered beneath the player and dynamic objects. Layers with z-index 1 and above are rendered above the player and dynamic objects.
- Efficient access and modification of tiles is required.
- Support for serialization/deserialization (saving/loading maps) in JSON format only (for document database and admin API).
- Algorithms for pathfinding, tile updates, and layer management.
- Each map can have links to other maps for each cardinal direction (north, south, east, west only), represented as the ID of the next map. Diagonal directions are not supported.

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

- `passable`: boolean (if false, tile is completely unpassable)
- `blocked_directions`: optional, array of objects, each with:
    - `direction`: enum (0 = north, 1 = east, 2 = south, 3 = west)
    - `block_inbound`: boolean
    - `block_outbound`: boolean
- `graphic_layers`: array of objects, each with:
    - `z_index`: integer
    - `graphic_id`: integer
    - `properties`: key/value pairs (map of strings)
- `warp`: optional, array [map_id, x, y] (all integers) to send the player to another location
- `trigger`: optional, string used for future scripting system
- `attributes`: key/value pairs (map of strings), intended to affect gameplay for the tile
