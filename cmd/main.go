package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"sync"

	"github.com/Odyssey-Classic/server/internal/server"

	"github.com/Odyssey-Classic/server/pb"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	// TODO: we don't ned this here, but not sure where we need it yet.
	var _ pb.GameMessage

	var registryURL string
	flag.StringVar(&registryURL, "registry", "http://local.fosteredgames.com:8080", "Registry URL")
	flag.Parse()

	cfg := server.Config{
		Ports: server.Ports{
			Admin:   GetUint16("ODY_ADMIN_PORT", 8081),
			Meta:    GetUint16("ODY_META_PORT", 8082),
			Network: GetUint16("ODY_NETWORK_PORT", 8080),
		},
		DataDir: GetString("ODY_DATA_DIR", "data"),
	}

	srv, err := server.NewServer(cfg,
		server.WithRegistry(registryURL),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	wg := new(sync.WaitGroup)
	srv.Start(ctx, wg)
	wg.Wait()
}
