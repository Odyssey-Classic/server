package network

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ClientMap map[*websocket.Conn]*Client

type Network struct {
	wg   *sync.WaitGroup
	port uint16
	once sync.Once

	clientGroup *sync.WaitGroup

	Out     chan any
	clients ClientMap
}

func New(port uint16) *Network {
	return &Network{
		port:        port,
		clientGroup: new(sync.WaitGroup),

		Out: make(chan any, 10),
	}
}

func (n *Network) Start(ctx context.Context, wg *sync.WaitGroup) error {
	var startErr error
	n.once.Do(func() {
		n.wg = wg
		n.wg.Add(1)
		go func() {
			startErr = n.start(ctx)
			n.wg.Done()
		}()
	})
	return startErr
}

func (n *Network) start(ctx context.Context) error {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", n.port),
	}
	server.Handler = n.wsConnect(ctx)
	server.BaseContext = func(listener net.Listener) context.Context { return ctx }

	go func() {
		<-ctx.Done()
		slog.Info("network shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		n.shutdown(shutdownCtx)
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("network shutdown error", "err", err)
		} else {
			slog.Info("network shutdown complete")
		}
	}()

	slog.Info("network API starting on :" + fmt.Sprintf("%d", n.port))
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("network server error", "err", err)
	}
	return err
}

func (n *Network) addClient(ctx context.Context, client *Client) {
	slog.Info("adding client", "remote addr", client.conn.RemoteAddr())
	n.clients[client.conn] = client
	n.Out <- client
	n.processClient(ctx, client)
}

func (n *Network) shutdown(_ context.Context) {
	slog.Info("shutting down clients")
	for _, client := range n.clients {
		err := client.close()
		if err != nil {
			slog.Error("error closing client", "error", err)
		}
	}
	slog.Info("clients shutdown")
}

// Stop shuts down the Network service
func (n *Network) Stop() {
	// Implement shutdown logic here
}
