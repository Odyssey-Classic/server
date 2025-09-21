```markdown
# Admin UI

The Admin UI is a desktop web application hosted via the server that allows administrators to manage their game server.

## Objective

Build a desktop-first web UI (served by the server) that supports administration workflows with a maps-first priority. The immediate scope is a maps editor MVP with the following constraints and goals:

- Provide a top-level menu for navigation and a landing/dashboard page showing server info and quick links.
- Primary focus: implement Map CRUD and an interactive map editor that supports single-tile edits.
- The editor must provide an explicit "Save" button that updates the server, and also autosave the in-progress map locally so closing and reopening the browser restores the user's work-in-progress.
  - Interface should make it clear the current map has edits pending upload to server.
- Map size is fixed to 17x17 tiles. Target browsers: Chrome and Firefox (desktop only).

## Technologies

- PixiJS for rendering maps and tilesets (used only for canvas rendering of maps/tiles).
- React for the UI shell, controls, and state management.
- SCSS for styling.

## Sections

These sections will be implemented later:
- Users
- Server Configuration
- Monsters
- Items and Equipment
- Live map browser

### Map Editing (detailed)

Requirements and behavior (maps-first MVP):

- Single-tile edits only for now. Implementation should remain extensible to add brushes/fills/multi-tile tools later.
- Central canvas renders the current 17x17 map using PixiJS. If no map is selected, show an empty state with actions to create or open a map.
- At each edge of the central canvas, render a preview strip showing two tiles deep from the linked map in that direction (north/south/east/west). This is a visual preview to help with edge continuity; preview tiles are not directly editable from the strip.
- Tile palette: a scrollable panel for selecting the active tile (visuals may be Pixi-rendered or an image grid until a tileset model exists).
- Explicit Save button: user-initiated save sends the current map to the server via the admin maps API.
- Local autosave: persist working copies to localStorage keyed by map id + user; on load, prefer the local draft if present and newer than the server version, and prompt the user to restore.
- Map links: maps may include links (north/east/south/west) referenced by id to support preview rendering.
- Desktop-only controls and pointer interactions; keyboard shortcuts can be added later.
- Undo/redo will be deferred to a future milestone.

Data model (suggested canonical JSON):

{
  "id": "map-1",
  "slug": "starting-village",
  "name": "Starting Village",
  "width": 17,
  "height": 17,
  "layers": [
    { "name": "ground", "data": [[0,1,2,...], ...] }
  ],
  "objects": [
    { "id": "obj-1", "type": "npc", "x": 5, "y": 10, "props": { "npcId": "greeting-npc" } }
  ],
  "links": { "north": "map-2", "east": "map-3" },
  "meta": { "createdBy": "admin", "createdAt": "2025-09-19T12:00:00Z" }
}

Server-side admin map operations are implemented in `../internal/services/admin/maps`.
Consult that package for exact method names and request/response shapes; the UI should call the server's admin maps service (via the server's existing HTTP handlers or adapters) rather than assuming specific REST endpoints.
UI layout sketch (maps editor):

- Top: global toolbar (app title, navigation menu, Save button, Undo/Redo placeholder, server status)
- Left: tile palette (select tile)
- Center: PixiJS canvas (map editor). Shows 17x17 tiles and 2-tile preview strips at edges.
- Right: inspector (map metadata, selected tile/object properties, link editor)
- Bottom: status bar (cursor coords, active tile id, local save indicator)

PixiJS + React integration recommendation:

- Implement a `MapCanvas` React component that mounts a Pixi.Application into a container ref on mount.
- Let Pixi handle rendering layers, tiles, and preview strips. React manages UI state (selected tile, working map JSON, edit mode).
- Expose a small imperative API from `MapCanvas` (e.g., loadMap(map), setActiveTile(tileId), applyTileEdit(x,y,tileId), panTo(x,y)) so other React components can command Pixi without forcing React re-renders on every tile edit.
- Keep the authoritative map model in React state or a light store (Context or Zustand). Batch updates to Pixi for rendering and to the React state only as needed for performance.

Local autosave strategy:

- On each tile edit, persist the working map to localStorage under key `admin-map-draft:{mapId}` with a timestamp.
- On load, check for a draft and compare timestamps with server data; offer the user to restore the draft or load the server copy.
- Explicit Save clears the local draft after a successful server update.

Acceptance criteria (maps-first MVP):

- Admin can list and open maps.
- Admin can create a new 17x17 map.
- Admin can select a tile from the palette and perform single-tile edits on the map.
- Editor shows two-tile preview strips for linked maps at edges.
- Local edits persist in localStorage and restore after browser close/reopen.
- Explicit Save updates the server map via the API.
- Desktop Chrome and Firefox supported.

Notes and future work:

- Add tileset metadata and a tileset data structure when available.
- Add server-side role-based UI controls later when auth/roles are implemented.
- Implement undo/redo, multi-tile tools, E2E tests, and performance optimizations in future iterations.
````
- API for making map edits exists in ../internal/services/admin/maps.
