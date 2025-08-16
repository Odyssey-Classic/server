package server

import (
	"context"
	"errors"
	"net/url"
	"sync"

	"github.com/Odyssey-Classic/server/internal/services/admin"
	"github.com/Odyssey-Classic/server/internal/services/game"
	"github.com/Odyssey-Classic/server/internal/services/meta"
	"github.com/Odyssey-Classic/server/internal/services/network"
)

type Server struct {
	once sync.Once

	admin   *admin.Admin
	meta    *meta.Meta
	network *network.Network
	game    *game.Game

	// Applied via Option
	registryURL *url.URL

	wg *sync.WaitGroup
}

func NewServer(cfg Config, options ...Option) (*Server, error) {
	server := &Server{
		wg: &sync.WaitGroup{},
	}

	server.admin = admin.New(cfg.Ports.Admin)
	server.meta = meta.New(cfg.Ports.Meta)
	server.network = network.New(cfg.Ports.Network)
	server.game = game.New(server.network.Out)

	// errors.Join will keep this value `nil` if no new errors are added.
	var optErrs error
	for _, opt := range options {
		optErrs = errors.Join(optErrs, opt(server))
	}

	return server, optErrs
}

func (s *Server) Start(ctx context.Context, wg *sync.WaitGroup) error {
	var err error
	s.once.Do(func() {
		err = s.start(ctx)
	})
	return err
}

func (s *Server) start(ctx context.Context) error {
	var startErr error
	startErr = errors.Join(startErr, s.admin.Start(ctx, s.wg))
	startErr = errors.Join(startErr, s.meta.Start(ctx, s.wg))
	startErr = errors.Join(startErr, s.network.Start(ctx, s.wg))
	startErr = errors.Join(startErr, s.game.Start(ctx, s.wg))

	if startErr != nil {
		s.stop()
		return startErr
	}

	s.wg.Wait()
	return nil
}

// stop shuts down all sub-services
func (s *Server) stop() {
	if s.admin != nil {
		s.admin.Stop()
	}
	if s.meta != nil {
		s.meta.Stop()
	}
	if s.network != nil {
		s.network.Stop()
	}
	if s.game != nil {
		s.game.Stop()
	}
}
