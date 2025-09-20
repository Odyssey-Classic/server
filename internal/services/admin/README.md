# Admin API

This package implements the Admin API for the Odyssey Classic server, focusing on map editing capabilities with an expandable structure for future admin features.

## Structure

```
internal/services/admin/api/
├── api.go              # Main API entry point and routing
├── utils/
│   └── responses.go    # Common response utilities
├── maps/
│   ├── api.go         # Maps CRUD operations
│   └── api_test.go    # Maps API tests
└── example/
    └── main.go        # Example server implementation
```

## Features

- **Modular Design**: Each admin feature (maps, users, settings, etc.) gets its own sub-package
- **Consistent Responses**: Standardized JSON responses using utility functions
- **RESTful Design**: Following REST principles for predictable API behavior
- **Chi Router**: Using Chi for its maintainability and modularity
- **Comprehensive Testing**: Unit tests for all endpoints

## Maps API Endpoints

| Endpoint           | Method | Description           |
|--------------------|--------|-----------------------|
| `/admin/maps`      | GET    | List all maps         |
| `/admin/maps`      | POST   | Create a new map      |
| `/admin/maps/{id}` | GET    | Get map details       |
| `/admin/maps/{id}` | PUT    | Update whole map      |
| `/admin/maps/{id}` | DELETE | Delete a map          |

## Usage

### Basic Server Setup

```go
package main

import (
    "net/http"
    "github.com/Odyssey-Classic/server/internal/services/admin/api"
)

func main() {
    adminAPI := api.New()
    
    server := &http.Server{
        Addr:    ":8080",
        Handler: adminAPI.Routes(),
    }
    
    server.ListenAndServe()
}
```

### Example Request/Response

**Create Map (POST /admin/maps)**

Request:
```json
{
  "name": "Forest of Dawn",
  "tags": ["forest", "beginner"],
  "attributes": {
    "difficulty": "easy",
    "theme": "nature"
  }
}
```

Response:
```json
{
  "id": 1,
  "name": "Forest of Dawn",
  "tags": ["forest", "beginner"],
  "attributes": {
    "difficulty": "easy",
    "theme": "nature"
  },
  "last_updated": "2025-08-27T12:00:00Z",
  "version": 1,
  "tiles": [...],
  "links": {}
}
```

## Data Types

The API uses the Map types defined in `/internal/game/maps`:

- `Map`: Complete map structure with tiles, metadata, and links
- `Tile`: Individual tile with graphics, passability, and attributes
- `MapLinks`: Links to adjacent maps by direction

## Future Expansion

The API structure is designed to easily accommodate new admin features:

```go
// In api.go setupRoutes()
a.router.Route("/admin", func(r chi.Router) {
    r.Mount("/maps", a.mapsAPI.Routes())
    r.Mount("/users", a.usersAPI.Routes())     // Future
    r.Mount("/settings", a.settingsAPI.Routes()) // Future
    r.Mount("/logs", a.logsAPI.Routes())       // Future
})
```

## Testing

Run tests for the maps API:

```bash
go test ./internal/services/admin/api/maps
```

Run all admin API tests:

```bash
go test ./internal/services/admin/api/...
```
