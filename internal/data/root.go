package data

import (
	"path/filepath"
)

// Root represents an abstraction over the server's data directory structure.
// It provides methods to access various subdirectories for different data types.
type Root interface {
	// MapsDir returns the path to the maps data directory
	MapsDir() string
}

// osRoot is an implementation of Root that uses the operating system's filesystem
type osRoot struct {
	baseDir string
}

// NewOSRoot creates a new Root implementation that uses the OS filesystem
// at the specified base directory path.
func NewOSRoot(baseDir string) Root {
	return &osRoot{
		baseDir: baseDir,
	}
}

// MapsDir returns the path to the maps subdirectory within the base data directory
func (r *osRoot) MapsDir() string {
	return filepath.Join(r.baseDir, "maps")
}
