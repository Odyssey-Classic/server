package maps

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	gamemaps "github.com/Odyssey-Classic/server/internal/game/maps"
	"github.com/Odyssey-Classic/server/internal/services/admin/maps/store"
	filestore "github.com/Odyssey-Classic/server/internal/services/admin/maps/store/file"
	"github.com/Odyssey-Classic/server/internal/services/admin/utils"
)

// API represents the maps admin API
type API struct {
	store store.MapStore
}

// NewFileBacked returns an API backed by the file store at root (default data/maps when empty).
func NewFileBacked(root string) *API {
	fs, err := filestore.New(root)
	if err != nil {
		panic(err)
	}
	return &API{store: fs}
}

// NewWithStore allows injecting a custom store (useful for alternate backends).
// Keep this constructor minimal and only if actually used elsewhere.
func NewWithStore(s store.MapStore) *API { return &API{store: s} }

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
	query := r.URL.Query().Get("q")
	maps, err := a.store.List(query)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to list maps")
		return
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
	name := strings.TrimSpace(mapData.Name)
	if name == "" {
		utils.WriteError(w, http.StatusBadRequest, "Name is required")
		return
	}
	m, err := a.store.Create(name)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create map")
		return
	}
	if err := utils.WriteJSON(w, http.StatusCreated, m); err != nil {
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

	m, err := a.store.Get(id)
	if err != nil {
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

	var mapData gamemaps.Map
	if err := json.NewDecoder(r.Body).Decode(&mapData); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Preserve the original ID
	mapData.ID = id
	if err := a.store.Update(&mapData); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update map")
		return
	}

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

	if err := a.store.Delete(id); err != nil {
		utils.WriteError(w, http.StatusNotFound, "Map not found")
		return
	}

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
