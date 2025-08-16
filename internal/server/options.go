package server

import "github.com/Odyssey-Classic/server/internal/services/registry"

// Go Options pattern to allow extensiblity

type Option func(*Server) error

// WithRegistry provides a url for an Odyssey Registry that the server MUST
// connect to or fail to start.
func WithRegistry(url string) Option {
	return func(s *Server) error {
		host, err := registry.ParseAndValidateURL(url)
		if err != nil {
			return err
		}
		s.registryURL = host
		return nil
	}
}
