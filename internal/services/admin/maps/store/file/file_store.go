package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	gamemaps "github.com/Odyssey-Classic/server/internal/game/maps"
)

// FileStore persists maps as JSON files under a root directory.
type FileStore struct {
	root   string
	mu     sync.RWMutex
	nextID int
}

// New creates a FileStore pointing at the provided root directory.
// If root is empty, it defaults to data/maps.
func New(root string) (*FileStore, error) {
	if root == "" {
		root = filepath.Join("data", "maps")
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}

	fs := &FileStore{root: root}
	// Initialize nextID by scanning existing files
	maxID, err := fs.scanMaxID()
	if err != nil {
		return nil, err
	}
	fs.nextID = maxID + 1
	if fs.nextID < 1 {
		fs.nextID = 1
	}
	return fs, nil
}

func (s *FileStore) scanMaxID() (int, error) {
	entries, err := os.ReadDir(s.root)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return 0, nil
		}
		return 0, err
	}
	maxID := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		// Expect filename like 000001.json
		base := strings.TrimSuffix(e.Name(), ".json")
		if id, err := strconv.Atoi(base); err == nil && id > maxID {
			maxID = id
		}
	}
	return maxID, nil
}

func (s *FileStore) pathFor(id int) string {
	return filepath.Join(s.root, fmt.Sprintf("%06d.json", id))
}

// Create a new map with only a name.
func (s *FileStore) Create(name string) (*gamemaps.Map, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++

	m := gamemaps.NewMap(id, name)
	// Ensure LastUpdated sensible
	m.LastUpdated = time.Now()

	if err := s.writeMap(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *FileStore) Get(id int) (*gamemaps.Map, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p := s.pathFor(id)
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var m gamemaps.Map
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *FileStore) Update(m *gamemaps.Map) error {
	if m == nil || m.ID <= 0 {
		return fmt.Errorf("invalid map or id")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	m.LastUpdated = time.Now()
	return s.writeMap(m)
}

func (s *FileStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.pathFor(id)
	return os.Remove(p)
}

func (s *FileStore) List(query string) ([]*gamemaps.Map, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.root)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return []*gamemaps.Map{}, nil
		}
		return nil, err
	}
	q := strings.ToLower(strings.TrimSpace(query))
	out := make([]*gamemaps.Map, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(s.root, e.Name()))
		if err != nil {
			continue
		}
		var m gamemaps.Map
		if err := json.Unmarshal(b, &m); err != nil {
			continue
		}
		if q == "" || strings.Contains(strings.ToLower(m.Name), q) {
			cp := m
			out = append(out, &cp)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (s *FileStore) writeMap(m *gamemaps.Map) error {
	p := s.pathFor(m.ID)
	// write to temp then rename for atomicity
	tmp, err := os.CreateTemp(s.root, "*.tmp")
	if err != nil {
		return err
	}
	enc := json.NewEncoder(tmp)
	enc.SetIndent("", "  ")
	if err := enc.Encode(m); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return os.Rename(tmp.Name(), p)
}
