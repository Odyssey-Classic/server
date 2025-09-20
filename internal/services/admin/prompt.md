
## Admin API Prompt: Map Editing

### Overview
Design an Admin API focused exclusively on map editing for the game server. The API should allow administrators to create, update, delete, and manage maps and their contents (such as tiles and objects). The design should be modular and easily expandable to support additional admin features in the future.

### Core Features (Map Editing)
- Create new maps
- List existing maps
- Update whole map data (name, description, size, tiles, objects, etc.)
- Delete maps

### Endpoints
| Endpoint           | Method | Description           |
|--------------------|--------|-----------------------|
| /admin/maps        | GET    | List all maps         |
| /admin/maps        | POST   | Create a new map      |
| /admin/maps/{id}   | GET    | Get map details       |
| /admin/maps/{id}   | PUT    | Update whole map      |
| /admin/maps/{id}   | DELETE | Delete a map          |

### Authentication & Authorization
Endpoints should be protected using JWT or similar mechanisms. Role-based access is recommended for future expansion.

### Tech Stack & Patterns
- Language: Go
- Framework: Chi (preferred for maintainability and modularity)
- RESTful design
- Modular structure for future features

### Example Requests & Responses

All request and response bodies for map operations should use the JSON structure defined by the Go types in `/internal/game/maps`. See the source code for `Map` and related types for the authoritative schema. Example:

> Refer to the Go types in `/internal/game/maps/map.go` and `/internal/game/maps/types.go` for the structure of map objects and their fields. API endpoints should serialize and deserialize map data according to these types.

### Additional Notes
- API should be versioned for future compatibility.
- All endpoints should return consistent error formats.
- Design should allow easy addition of new admin features (e.g., user management, settings) in the future.
