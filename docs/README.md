# Odyssey Server

The Server consists of 4 parts
- Game Logic
- Networking
- Meta
- Admin

Game Logic is dedicated to maintaining the game simulation.  
It handles resolving player actions, game events, and sharing results.  
Data persistence is handled by this layer.

Networking handles the dedicated, real time, bidirectional communication between
clients and the server.

Meta handles providing asynchronous data and actions with the server.  
e.g. listing a player's characters for selection before connecting through
the Networking layer.

Admin handles administrative tasks.  
e.g. an API that allows admins to upload Map updates.

## Admin UI hosting

The server hosts the Admin UI (React + Vite) directly from the Admin service.

- Build output: The UI is built into `internal/web/dist` and embedded in the Go binary. See the Vite config at `ui/vite.config.ts`.
- Embedding and serving: Assets under `internal/web/dist` are embedded and served as a Single Page App (SPA). See `internal/web/spa.go` for the embed and handler, and `internal/services/admin/admin.go` for router wiring.
- Routes:
	- UI is served at the Admin service root path (`/`).
	- Static assets are served under `/assets/...`.
	- Existing Admin API routes remain under `/admin/...`.
- Build pipeline: The standard server build runs the UI build first and then compiles the Go binary, so no extra step is required. See `Makefile`.
- Development flow: For UI development, use the existing dev workflow that runs the server and the Vite dev server concurrently. The production build is embedded and served by the server.

## Player Join Flow
Client: sends request to Meta with Id.  
Meta: returns character list, any other info needed to start playing.  
Client: requests connection with Id and data returned from meta.  
Networking: upgrades to websocket persistent connection.  
Game: puts character in world.  
Client <-> Netowrk <-> Game, Player playing.  

This way, Player objects are always full players and not in  
"connected but not playing" states.
