package maps

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	gamemaps "github.com/Odyssey-Classic/server/internal/game/maps"
	"github.com/Odyssey-Classic/server/internal/services/admin/utils"
)

// API represents the maps admin API
type API struct {
	// In a real implementation, this would be a database or service layer
	// For now, we'll use an in-memory store as a placeholder
	maps   map[int]*gamemaps.Map
	nextID int
}

// New creates a new maps API instance
func New() *API {
	return &API{
		maps:   make(map[int]*gamemaps.Map),
		nextID: 1,
	}
}

// Routes returns the chi router for maps endpoints
func (a *API) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", a.listMaps)
	r.Post("/", a.createMap)
	r.Get("/{id}", a.getMap)
	r.Put("/{id}", a.updateMap)
	r.Delete("/{id}", a.deleteMap)

	return r
}

// listMaps handles GET /admin/maps - List all maps
func (a *API) listMaps(w http.ResponseWriter, r *http.Request) {
	maps := make([]*gamemaps.Map, 0, len(a.maps))
	for _, m := range a.maps {
		maps = append(maps, m)
	}

	if err := utils.WriteJSON(w, http.StatusOK, maps); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

// createMap handles POST /admin/maps - Create a new map
func (a *API) createMap(w http.ResponseWriter, r *http.Request) {
	var mapData gamemaps.Map
	if err := json.NewDecoder(r.Body).Decode(&mapData); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Assign ID and store
	mapData.ID = a.nextID
	a.nextID++

	// Create new map with proper initialization
	newMap := gamemaps.NewMap(mapData.ID, mapData.Name)
	newMap.Tags = mapData.Tags
	newMap.Attributes = mapData.Attributes
	newMap.Tiles = mapData.Tiles
	newMap.Links = mapData.Links

	a.maps[newMap.ID] = newMap

	if err := utils.WriteJSON(w, http.StatusCreated, newMap); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

// getMap handles GET /admin/maps/{id} - Get map details
func (a *API) getMap(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid map ID")
		return
	}

	m, exists := a.maps[id]
	if !exists {
		utils.WriteError(w, http.StatusNotFound, "Map not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, m); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

// updateMap handles PUT /admin/maps/{id} - Update whole map
func (a *API) updateMap(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid map ID")
		return
	}

	_, exists := a.maps[id]
	if !exists {
		utils.WriteError(w, http.StatusNotFound, "Map not found")
		return
	}

	var mapData gamemaps.Map
	if err := json.NewDecoder(r.Body).Decode(&mapData); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Preserve the original ID
	mapData.ID = id
	a.maps[id] = &mapData

	response := map[string]interface{}{
		"success":        true,
		"updated_map_id": id,
	}

	if err := utils.WriteJSON(w, http.StatusOK, response); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

// deleteMap handles DELETE /admin/maps/{id} - Delete a map
func (a *API) deleteMap(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid map ID")
		return
	}

	_, exists := a.maps[id]
	if !exists {
		utils.WriteError(w, http.StatusNotFound, "Map not found")
		return
	}

	delete(a.maps, id)

	response := map[string]interface{}{
		"success":    true,
		"deleted_id": id,
		"message":    fmt.Sprintf("Map %d deleted successfully", id),
	}

	if err := utils.WriteJSON(w, http.StatusOK, response); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}
