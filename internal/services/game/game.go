package game

import (
	"context"
	"log/slog"
	"sync"
)

type Game struct {
	wg   *sync.WaitGroup
	once sync.Once

	network chan any
}

func New(network chan any) *Game {
	return &Game{
		network: network,
	}
}

func (g *Game) Start(ctx context.Context, wg *sync.WaitGroup) error {
	var startErr error
	g.once.Do(func() {
		g.wg = wg
		g.wg.Add(1)
		go func() {
			startErr = g.start(ctx)
			g.wg.Done()
			slog.Info("game shutdown complete")
		}()
	})
	return startErr
}

func (g *Game) start(ctx context.Context) error {
	for {
		select {
		case msg := <-g.network:
			slog.Info("game received message: ", "message", msg)
			// msg.Type = new Client
			//   game.NewPlayer(client)
		case <-ctx.Done():
			slog.Info("game shutting down")
			return nil
		}
	}
}

// Stop shuts down the Game service
func (g *Game) Stop() {
	// Implement shutdown logic here
}
