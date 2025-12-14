package data

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RootTestSuite struct {
	suite.Suite
}

func TestRootTestSuite(t *testing.T) {
	suite.Run(t, new(RootTestSuite))
}

func (s *RootTestSuite) TestNewOSRoot() {
	baseDir := "/test/data"
	root := NewOSRoot(baseDir)

	s.NotNil(root, "NewOSRoot should return a non-nil Root")
}

func (s *RootTestSuite) TestMapsDir() {
	baseDir := "/test/data"
	root := NewOSRoot(baseDir)

	mapsDir := root.MapsDir()
	expected := filepath.Join(baseDir, "maps")

	s.Equal(expected, mapsDir, "MapsDir should return the correct path")
}

func (s *RootTestSuite) TestMapsDir_RelativePath() {
	baseDir := "data"
	root := NewOSRoot(baseDir)

	mapsDir := root.MapsDir()
	expected := filepath.Join("data", "maps")

	s.Equal(expected, mapsDir, "MapsDir should work with relative paths")
}
