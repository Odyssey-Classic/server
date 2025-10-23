package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Odyssey-Classic/server/internal/data"
	"github.com/stretchr/testify/suite"
)

// AdminAPITestSuite defines the test suite for Admin API integration tests
type AdminAPITestSuite struct {
	suite.Suite
	api *API
}

// SetupTest runs before each test method
func (s *AdminAPITestSuite) SetupTest() {
	// Use a per-test temporary data directory via data.Root abstraction
	tmp := s.T().TempDir()
	s.api = api(data.NewOSRoot(tmp))
}

// TestMiddlewareSetup tests that the API sets up middleware correctly
func (s *AdminAPITestSuite) TestMiddlewareSetup() {
	// Test with a valid route to see middleware working
	req := httptest.NewRequest(http.MethodGet, "/admin/maps", nil)
	w := httptest.NewRecorder()

	s.api.ServeHTTP(w, req)

	// Check that Content-Type middleware is applied for valid requests
	s.Equal("application/json", w.Header().Get("Content-Type"))

	// Should get 200 OK for the maps endpoint
	s.Equal(http.StatusOK, w.Code)
}

// TestRoutesSetup tests that routes are properly mounted
func (s *AdminAPITestSuite) TestRoutesSetup() {
	// Test that maps routes are mounted under /admin/maps
	req := httptest.NewRequest(http.MethodGet, "/admin/maps", nil)
	w := httptest.NewRecorder()

	s.api.ServeHTTP(w, req)

	// Should get 200 OK (empty list) rather than 404, indicating route is mounted
	s.Equal(http.StatusOK, w.Code)
}

// TestAPIStructure tests the API structure and composition
func (s *AdminAPITestSuite) TestAPIStructure() {
	// Verify API is properly constructed
	s.NotNil(s.api, "API should be created")
	s.NotNil(s.api.Routes(), "API should have routes")

	// Test that the API implements http.Handler
	var handler http.Handler = s.api
	s.NotNil(handler, "API should implement http.Handler")
}

// TestHTTPHandlerInterface tests that the API properly implements http.Handler
func (s *AdminAPITestSuite) TestHTTPHandlerInterface() {
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()

	// This should not panic and should handle the request
	s.NotPanics(func() {
		s.api.ServeHTTP(w, req)
	}, "ServeHTTP should not panic")

	// Should return some response (even if 404)
	s.True(w.Code >= 200 && w.Code < 600, "Should return valid HTTP status code")
}

// TestAdminAPI_Integration runs the complete test suite
func TestAdminAPI_Integration(t *testing.T) {
	suite.Run(t, new(AdminAPITestSuite))
}
