package maps

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"

	gamemaps "github.com/Odyssey-Classic/server/internal/game/maps"
)

// MapsAPITestSuite defines the test suite for Maps API tests
type MapsAPITestSuite struct {
	suite.Suite
	api    *API
	router chi.Router
}

// SetupTest runs before each test method
func (s *MapsAPITestSuite) SetupTest() {
	s.api = New()
	s.router = chi.NewRouter()
	s.router.Mount("/admin/maps", s.api.Routes())
}

// TestListMaps_Empty tests listing maps when no maps exist
func (s *MapsAPITestSuite) TestListMaps_Empty() {
	req := httptest.NewRequest(http.MethodGet, "/admin/maps", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var maps []*gamemaps.Map
	err := json.NewDecoder(w.Body).Decode(&maps)
	s.NoError(err, "Failed to decode response")

	s.Len(maps, 0, "Expected 0 maps")
}

// TestCreateMap tests creating a new map
func (s *MapsAPITestSuite) TestCreateMap() {
	newMap := gamemaps.Map{
		Name: "Test Map",
		Tags: []string{"test", "example"},
	}

	body, _ := json.Marshal(newMap)
	req := httptest.NewRequest(http.MethodPost, "/admin/maps", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)

	var createdMap gamemaps.Map
	err := json.NewDecoder(w.Body).Decode(&createdMap)
	s.NoError(err, "Failed to decode response")

	s.Equal(1, int(createdMap.ID))
	s.Equal("Test Map", createdMap.Name)
}

// TestGetMap tests retrieving an existing map
func (s *MapsAPITestSuite) TestGetMap() {
	// First create a map
	newMap := gamemaps.Map{
		Name: "Test Map",
		Tags: []string{"test", "example"},
	}

	body, _ := json.Marshal(newMap)
	req := httptest.NewRequest(http.MethodPost, "/admin/maps", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusCreated, w.Code, "Map creation should succeed")

	// Now retrieve it
	req = httptest.NewRequest(http.MethodGet, "/admin/maps/1", nil)
	w = httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var retrievedMap gamemaps.Map
	err := json.NewDecoder(w.Body).Decode(&retrievedMap)
	s.NoError(err, "Failed to decode response")

	s.Equal("Test Map", retrievedMap.Name)
}

// TestGetMap_NotFound tests retrieving a non-existent map
func (s *MapsAPITestSuite) TestGetMap_NotFound() {
	req := httptest.NewRequest(http.MethodGet, "/admin/maps/999", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusNotFound, w.Code)
}

// TestUpdateMap tests updating an existing map
func (s *MapsAPITestSuite) TestUpdateMap() {
	// First create a map
	newMap := gamemaps.Map{
		Name: "Test Map",
		Tags: []string{"test", "example"},
	}

	body, _ := json.Marshal(newMap)
	req := httptest.NewRequest(http.MethodPost, "/admin/maps", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusCreated, w.Code, "Map creation should succeed")

	// Now update it
	updatedMap := gamemaps.Map{
		Name: "Updated Test Map",
		Tags: []string{"updated", "test"},
	}

	body, _ = json.Marshal(updatedMap)
	req = httptest.NewRequest(http.MethodPut, "/admin/maps/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	s.NoError(err, "Failed to decode response")

	success, ok := response["success"].(bool)
	s.True(ok, "Response should contain success field")
	s.True(success, "Expected success to be true")
}

// TestDeleteMap tests deleting an existing map
func (s *MapsAPITestSuite) TestDeleteMap() {
	// First create a map
	newMap := gamemaps.Map{
		Name: "Test Map",
		Tags: []string{"test", "example"},
	}

	body, _ := json.Marshal(newMap)
	req := httptest.NewRequest(http.MethodPost, "/admin/maps", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusCreated, w.Code, "Map creation should succeed")

	// Now delete it
	req = httptest.NewRequest(http.MethodDelete, "/admin/maps/1", nil)
	w = httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	s.NoError(err, "Failed to decode response")

	success, ok := response["success"].(bool)
	s.True(ok, "Response should contain success field")
	s.True(success, "Expected success to be true")
}

// TestDeleteMap_NotFound tests deleting a non-existent map
func (s *MapsAPITestSuite) TestDeleteMap_NotFound() {
	req := httptest.NewRequest(http.MethodDelete, "/admin/maps/999", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusNotFound, w.Code)
}

// TestMapsAPI runs the complete test suite
func TestMapsAPI(t *testing.T) {
	suite.Run(t, new(MapsAPITestSuite))
}
